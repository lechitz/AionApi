package handler

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

// GetByID is the resolver for the tagByID field.
func (h *Handler) GetByID(_ context.Context, _ string) (*model.Tag, error) {
	// got, err := h.TagService.GetByID(ctx, id)
	// if err != nil {
	//	return nil, err
	// }
	// return toModelOut(got), nil
	return nil, errors.New("not implemented")
}
