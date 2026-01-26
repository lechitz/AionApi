package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetLatest retrieves the N most recent chat history entries for a user from cache.
func (s *Store) GetLatest(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameChatHistoryGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.Int(AttributeLimit, limit),
	))
	defer span.End()

	s.logger.InfowCtx(ctx, LogGetHistory,
		commonkeys.UserID, userID,
		AttributeLimit, limit,
	)

	key := buildKey(userID)
	data, err := s.cache.Get(ctx, key)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrGetHistory)
		s.logger.ErrorwCtx(ctx, ErrGetHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, fmt.Errorf("%s: %w", ErrGetHistory, err)
	}

	// Cache miss - return empty slice.
	if data == "" {
		span.SetAttributes(attribute.Bool(AttributeCacheHit, false))
		span.SetStatus(codes.Ok, LogHistoryCacheMiss)
		s.logger.InfowCtx(ctx, LogHistoryCacheMiss, commonkeys.UserID, userID)
		return []domain.ChatHistory{}, nil
	}

	var history []domain.ChatHistory
	if err := json.Unmarshal([]byte(data), &history); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrUnmarshalHistory)
		s.logger.ErrorwCtx(ctx, ErrUnmarshalHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return nil, fmt.Errorf("%s: %w", ErrUnmarshalHistory, err)
	}

	// Apply limit and return most recent.
	if len(history) > limit {
		history = history[:limit]
	}

	span.SetAttributes(
		attribute.Bool(AttributeCacheHit, true),
		attribute.Int(AttributeResultsCount, len(history)),
	)
	span.SetStatus(codes.Ok, LogHistoryRetrieved)

	s.logger.InfowCtx(ctx, LogHistoryRetrieved,
		commonkeys.UserID, userID,
		AttributeResultsCount, len(history),
	)

	return history, nil
}
