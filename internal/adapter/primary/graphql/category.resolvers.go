package graphql

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// --- Mutations ---

func (m *mutationResolver) CreateCategory(ctx context.Context, input model.CreateCategoryInput) (*model.Category, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.CategoryController().Create(ctx, input, uid)
}

func (m *mutationResolver) UpdateCategory(ctx context.Context, input model.UpdateCategoryInput) (*model.Category, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return m.CategoryController().Update(ctx, input, uid)
}

func (m *mutationResolver) SoftDeleteCategory(ctx context.Context, input model.DeleteCategoryInput) (bool, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	id, err := strconv.ParseUint(input.ID, 10, 64)
	if err != nil {
		return false, err
	}
	if err := m.CategoryController().SoftDelete(ctx, id, uid); err != nil {
		return false, err
	}
	return true, nil
}

// --- Queries ---

func (q *queryResolver) Categories(ctx context.Context) ([]*model.Category, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.CategoryController().ListAll(ctx, uid)
}

func (q *queryResolver) CategoryByID(ctx context.Context, id string) (*model.Category, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	cid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	return q.CategoryController().GetByID(ctx, cid, uid)
}

func (q *queryResolver) CategoryByName(ctx context.Context, name string) (*model.Category, error) {
	uid, _ := ctx.Value(ctxkeys.UserID).(uint64)
	return q.CategoryController().GetByName(ctx, name, uid)
}

// Ensure Resolver implements interfaces
var _ MutationResolver = (*mutationResolver)(nil)
var _ QueryResolver = (*queryResolver)(nil)
