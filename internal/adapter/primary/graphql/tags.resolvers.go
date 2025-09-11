package graphql

import (
	"context"
	"fmt"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

// CreateTag is the resolver for the createTag field.
func (m *mutationResolver) CreateTag(ctx context.Context, input model.CreateTagInput) (*model.Tag, error) {
	return nil, fmt.Errorf("not implemented: CreateTag")
}

// Tags is the resolver for the tag field.
func (q *queryResolver) Tags(ctx context.Context) ([]*model.Tag, error) {
	return nil, fmt.Errorf("not implemented: Tags")
}

// TagByID is the resolver for the tagById field.
func (q *queryResolver) TagByID(ctx context.Context, id string) (*model.Tag, error) {
	return nil, fmt.Errorf("not implemented: TagByID")
}
