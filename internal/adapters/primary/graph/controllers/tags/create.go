package tags

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
)

func (h *Handler) Create(ctx context.Context, in model.CreateTagInput) (*model.Tag, error) {
	//td, err := toDomainCreate(in)
	//if err != nil {
	//	return nil, err
	//}
	//
	//created, err := h.TagService.Create(ctx, td)
	//if err != nil {
	//	return nil, err
	//}
	//return toModelOut(created), nil
	return nil, nil
}
