package usecase

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// SaveChatHistory persists a chat interaction to the database for future reference.
// It creates a new ChatHistory entry with the user's message, AI's response, and metadata.
func (s *ChatService) SaveChatHistory(ctx context.Context, userID uint64, message, response string, tokensUsed int, functionCalls map[string]string) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSaveChatHistory)
	defer span.End()

	span.SetAttributes(
		attribute.String(AttrUserID, strconv.FormatUint(userID, 10)),
		attribute.Int("message_length", len(message)),
		attribute.Int("response_length", len(response)),
		attribute.Int(AttrTokensUsed, tokensUsed),
	)

	s.logger.InfowCtx(ctx, LogSavingChatHistory,
		LogKeyUserID, userID,
		LogKeyMessageLength, len(message),
		LogKeyResponseLength, len(response),
		LogKeyTokensUsed, tokensUsed,
	)

	// Create domain model
	chatHistory := domain.ChatHistory{
		UserID:        userID,
		Message:       message,
		Response:      response,
		TokensUsed:    tokensUsed,
		FunctionCalls: functionCalls,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to repository (persistent storage)
	saved, err := s.chatHistoryRepo.Save(ctx, chatHistory)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedToSaveChatHistory)
		s.logger.ErrorwCtx(ctx, LogFailedToSaveChatHistory,
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		return err
	}

	// Also save to cache (fast access for conversation context)
	// Don't fail if cache update fails - it's not critical
	if err := s.chatHistoryCache.Add(ctx, userID, saved); err != nil {
		s.logger.WarnwCtx(ctx, "Failed to update chat history cache (non-blocking)",
			LogKeyError, err.Error(),
			LogKeyUserID, userID,
		)
		// Don't return error - cache miss is acceptable
	}

	span.SetAttributes(
		attribute.String("chat_id", strconv.FormatUint(saved.ChatID, 10)),
	)
	span.SetStatus(codes.Ok, StatusChatHistorySaved)

	s.logger.InfowCtx(ctx, LogChatHistorySaved,
		LogKeyUserID, userID,
		"chat_id", saved.ChatID,
		LogKeyTokensUsed, saved.TokensUsed,
	)

	return nil
}
