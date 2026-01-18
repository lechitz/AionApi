package graphql

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// ChatHistory is the resolver for the chatHistory field.
func (r *queryResolver) ChatHistory(ctx context.Context, limit *int32, offset *int32) ([]*model.ChatMessage, error) {
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)

	// Default values
	limitVal := 10
	offsetVal := 0

	if limit != nil {
		limitVal = int(*limit)
	}
	if offset != nil {
		offsetVal = int(*offset)
	}

	return r.ChatController().GetChatHistory(ctx, userID, limitVal, offsetVal)
}

// ChatContext is the resolver for the chatContext field.
func (r *queryResolver) ChatContext(ctx context.Context) (*model.ChatContext, error) {
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return r.ChatController().GetChatContext(ctx, userID)
}
