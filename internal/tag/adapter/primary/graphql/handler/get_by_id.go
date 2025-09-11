package handler

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

func (h *Handler) GetByID(ctx context.Context, id string) (*model.Tag, error) {
	//got, err := h.TagService.GetByID(ctx, id)
	//if err != nil {
	//	return nil, err
	//}
	//return toModelOut(got), nil
	return nil, nil
}
