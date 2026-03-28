package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/chat/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetChatHistory retrieves chat conversation history for a user with pagination support.
// Returns the most recent conversations first (ordered by created_at DESC).
func (s *ChatService) GetChatHistory(ctx context.Context, userID uint64, limit, offset int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetChatHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
		attribute.Int("offset", offset),
	)

	s.logger.InfowCtx(ctx, LogGettingChatHistory,
		LogKeyUserID, userID,
		"limit", limit,
		"offset", offset,
	)

	// Get from repository
	histories, err := s.chatHistoryRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedToGetChatHistory)
		s.logger.ErrorwCtx(ctx, LogFailedToGetChatHistory,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("count", len(histories)),
	)
	span.SetStatus(codes.Ok, StatusChatHistoryRetrieved)

	s.logger.InfowCtx(ctx, LogChatHistoryRetrieved,
		LogKeyUserID, userID,
		"count", len(histories),
	)

	return histories, nil
}

// GetLatestChatHistory retrieves the N most recent chat entries for a user.
// This is a convenience method for getting recent context without pagination.
func (s *ChatService) GetLatestChatHistory(ctx context.Context, userID uint64, limit int) ([]domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetChatHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
		attribute.Int("limit", limit),
	)

	s.logger.InfowCtx(ctx, LogGettingChatHistory,
		LogKeyUserID, userID,
		"limit", limit,
	)

	// Get latest from repository
	histories, err := s.chatHistoryRepo.GetLatest(ctx, userID, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedToGetChatHistory)
		s.logger.ErrorwCtx(ctx, LogFailedToGetChatHistory,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return nil, err
	}

	span.SetAttributes(
		attribute.Int("count", len(histories)),
	)
	span.SetStatus(codes.Ok, StatusChatHistoryRetrieved)

	s.logger.InfowCtx(ctx, LogChatHistoryRetrieved,
		LogKeyUserID, userID,
		"count", len(histories),
	)

	return histories, nil
}
