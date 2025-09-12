package handler

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

// Create is the resolver for the createTag field.
func (h *Handler) Create(_ context.Context, _ model.CreateTagInput) (*model.Tag, error) {
	// td, err := toDomainCreate(in)
	// if err != nil {
	//	return nil, err
	// }
	//
	// created, err := h.TagService.Create(ctx, td)
	// if err != nil {
	//	return nil, err
	// }
	// return toModelOut(created), nil
	return nil, errors.New("not implemented")
}
