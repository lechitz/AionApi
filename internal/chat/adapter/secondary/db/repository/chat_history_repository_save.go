package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/chat/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/chat/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Save inserts a new chat history entry and returns it, or an error on failure.
func (r *ChatHistoryRepository) Save(ctx context.Context, chatHistory domain.ChatHistory) (domain.ChatHistory, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSaveRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(chatHistory.UserID, 10)),
		attribute.String(commonkeys.Operation, OpSave),
		attribute.Int("message_length", len(chatHistory.Message)),
		attribute.Int("response_length", len(chatHistory.Response)),
	))
	defer span.End()

	row := mapper.ChatHistoryToDB(chatHistory)

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpSave)
		r.logger.ErrorwCtx(ctx, ErrSaveChatMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(chatHistory.UserID, 10),
		)
		return domain.ChatHistory{}, fmt.Errorf("save chat history: %w", err)
	}

	saved := mapper.ChatHistoryFromDB(row)

	span.SetAttributes(
		attribute.String("chat_id", strconv.FormatUint(saved.ChatID, 10)),
		attribute.Int("tokens_used", saved.TokensUsed),
	)
	span.SetStatus(codes.Ok, StatusChatSaved)

	r.logger.InfowCtx(ctx, StatusChatSaved,
		commonkeys.UserID, strconv.FormatUint(saved.UserID, 10),
		"chat_id", saved.ChatID,
		"tokens_used", saved.TokensUsed,
	)

	return saved, nil
}
