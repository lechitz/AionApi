package usecase

import (
	"context"
	"strconv"

	"github.com/lechitz/aion-api/internal/chat/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetChatContext retrieves aggregated context data for AI to generate better responses.
// This includes recent records, available categories/tags, and recent chat history.
func (s *ChatService) GetChatContext(ctx context.Context, userID uint64) (*domain.ChatContext, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetChatContext)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
	)

	s.logger.InfowCtx(ctx, LogGettingChatContext, LogKeyUserID, userID)

	chatContext := &domain.ChatContext{}

	// Get recent chat history (non-blocking errors)
	if histories, err := s.chatHistoryRepo.GetLatest(ctx, userID, 5); err == nil {
		chatContext.RecentChats = histories
	} else {
		s.logger.WarnwCtx(ctx, "Failed to get recent chats for context (non-critical)",
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
	}

	// Note: Categories, Tags, and Records would be retrieved here if we had access to those services.
	// For now, we focus on chat history since we have that repository.
	// TODO: Inject CategoryService, TagService, RecordService to get complete context

	span.SetAttributes(
		attribute.Int("recent_chats_count", len(chatContext.RecentChats)),
	)
	span.SetStatus(codes.Ok, StatusChatContextRetrieved)

	s.logger.InfowCtx(ctx, LogChatContextRetrieved,
		LogKeyUserID, userID,
		"recent_chats", len(chatContext.RecentChats),
	)

	return chatContext, nil
}
