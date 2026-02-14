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

// ChatDataPack is the resolver for the chatDataPack field.
func (r *queryResolver) ChatDataPack(ctx context.Context, limitRecords *int32, includeStats bool) (*model.ChatDataPack, error) {
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)

	lim := 10
	if limitRecords != nil && *limitRecords > 0 {
		lim = int(*limitRecords)
	}

	categories, err := r.CategoryController().ListAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	tags, err := r.TagController().GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	recentRecords, err := r.RecordController().ListLatest(ctx, userID, lim)
	if err != nil {
		return nil, err
	}

	result := &model.ChatDataPack{
		Categories:    categories,
		Tags:          tags,
		RecentRecords: recentRecords,
	}

	if includeStats {
		stats, statsErr := r.UserStats(ctx)
		if statsErr != nil {
			return nil, statsErr
		}
		result.UserStats = stats
	}

	return result, nil
}
