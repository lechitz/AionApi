package handler

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
)

func (h *Handler) ListAll(ctx context.Context) ([]*model.Tag, error) {
	//all, err := h.TagService.ListAll(ctx)
	//if err != nil {
	//	return nil, err
	//}
	//out := make([]*model.Tag, len(all))
	//for i, t := range all {
	//	out[i] = toModelOut(t)
	//}
	//return out, nil
	return nil, nil
}
