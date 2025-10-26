package graphql

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// CreateTag is the resolver for the createTag field.
func (m *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.TagController().Create(ctx, input, uid)
}

// TagByName is the resolve for the tagByName field.
func (q *queryResolver) TagByName(ctx context.Context, tagName string) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.TagController().GetByName(ctx, tagName, uid)
}

// TagByID is the resolve for the tagByID field.
func (q *queryResolver) TagByID(ctx context.Context, tagID string) (*model.Tag, error) {
	id, err := strconv.ParseUint(tagID, 10, 64)
	if err != nil {
		return nil, err
	}
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.TagController().GetByID(ctx, id, uid)
}

func (q *queryResolver) TagByCategoryID(ctx context.Context, categoryID string) ([]*model.Tag, error) {
	id, err := strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		return nil, err
	}
	userID, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.TagController().GetByCategoryID(ctx, id, userID)
}

// GetAll is the resolver for the allTags field.
func (q *queryResolver) GetAll(ctx context.Context) ([]*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.TagController().GetAll(ctx, uid)
}
