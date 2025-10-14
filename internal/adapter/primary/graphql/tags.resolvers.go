package graphql

import (
	"context"
	"errors"

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
func (q *queryResolver) TagByID(ctx context.Context, tagID uint64) (*model.Tag, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.TagController().GetByID(ctx, tagID, uid)
}

// Tags is the resolver for the tag field.
func (q *queryResolver) Tags(_ context.Context) ([]*model.Tag, error) {
	return nil, errors.New("not implemented: Tags")
}
