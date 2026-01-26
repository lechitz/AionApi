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

// Add adds a new chat history entry to cache (prepends to the list).
func (s *Store) Add(ctx context.Context, userID uint64, history domain.ChatHistory) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameChatHistorySave, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	s.logger.InfowCtx(ctx, LogSetHistory, commonkeys.UserID, userID)

	key := buildKey(userID)

	// Get existing history.
	existingHistory, err := s.GetLatest(ctx, userID, historyLimitDefault)
	if err != nil {
		// If error, start with empty history.
		existingHistory = []domain.ChatHistory{}
	}

	// Prepend new entry (most recent first).
	newHistory := append([]domain.ChatHistory{history}, existingHistory...)

	// Limit to most recent entries.
	if len(newHistory) > historyLimitDefault {
		newHistory = newHistory[:historyLimitDefault]
	}

	// Marshal to JSON.
	data, err := json.Marshal(newHistory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrMarshalHistory)
		s.logger.ErrorwCtx(ctx, ErrMarshalHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return fmt.Errorf("%s: %w", ErrMarshalHistory, err)
	}

	// Set in cache with TTL.
	if err := s.cache.Set(ctx, key, string(data), defaultTTL); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrSetHistory)
		s.logger.ErrorwCtx(ctx, ErrSetHistory,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, userID,
		)
		return fmt.Errorf("%s: %w", ErrSetHistory, err)
	}

	span.SetAttributes(attribute.Int(AttributeHistorySize, len(newHistory)))
	span.SetStatus(codes.Ok, LogHistorySet)

	s.logger.InfowCtx(ctx, LogHistorySet,
		commonkeys.UserID, userID,
		AttributeHistorySize, len(newHistory),
	)

	return nil
}
