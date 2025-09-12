package graphql

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

// CreateTag is the resolver for the createTag field.
func (m *mutationResolver) CreateTag(_ context.Context, _ model.CreateTagInput) (*model.Tag, error) {
	return nil, errors.New("not implemented: CreateTag")
}

// Tags is the resolver for the tag field.
func (q *queryResolver) Tags(_ context.Context) ([]*model.Tag, error) {
	return nil, errors.New("not implemented: Tags")
}

// TagByID is the resolver for the tagById field.
func (q *queryResolver) TagByID(_ context.Context, _ string) (*model.Tag, error) {
	return nil, errors.New("not implemented: TagByID")
}
