package handler

import (
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/core/domain"
)

func toDomainCreate(in model.CreateCategoryInput, userID uint64) domain.Category {
	c := domain.Category{UserID: userID, Name: in.Name}
	if in.Description != nil {
		c.Description = *in.Description
	}
	if in.ColorHex != nil {
		c.Color = *in.ColorHex
	}
	if in.Icon != nil {
		c.Icon = *in.Icon
	}
	return c
}

// Ajuste esta função conforme seu schema real de UpdateCategoryInput.
// Assumimos campos opcionais: id, name, description, colorHex, icon.
func toDomainUpdate(in model.UpdateCategoryInput, userID uint64) domain.Category {
	var id uint64
	if in.ID != nil {
		id, _ = strconv.ParseUint(*in.ID, 10, 64)
	}
	c := domain.Category{
		ID:     id,
		UserID: userID,
	}
	if in.Name != nil {
		c.Name = *in.Name
	}
	if in.Description != nil {
		c.Description = *in.Description
	}
	if in.ColorHex != nil {
		c.Color = *in.ColorHex
	}
	if in.Icon != nil {
		c.Icon = *in.Icon
	}
	return c
}

func toModel(d domain.Category) *model.Category {
	out := &model.Category{
		ID:     strconv.FormatUint(d.ID, 10),
		UserID: strconv.FormatUint(d.UserID, 10),
		Name:   d.Name,
	}
	if d.Description != "" {
		out.Description = &d.Description
	}
	if d.Color != "" {
		out.ColorHex = &d.Color
	}
	if d.Icon != "" {
		out.Icon = &d.Icon
	}
	return out
}

func toModelList(dd []domain.Category) []*model.Category {
	if len(dd) == 0 {
		return []*model.Category{}
	}
	out := make([]*model.Category, 0, len(dd))
	for _, d := range dd {
		out = append(out, toModel(d))
	}
	return out
}
