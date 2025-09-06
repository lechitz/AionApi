package handler

import (
	"strconv"

	"github.com/lechitz/AionApi/internal/adapter/server/graph/model"
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

func toDomainUpdate(in model.UpdateCategoryInput, userID uint64) (domain.Category, error) {
	id, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return domain.Category{}, err
	}
	c := domain.Category{ID: id, UserID: userID}
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
	return c, nil
}

func toModelOut(c domain.Category) *model.Category {
	return &model.Category{
		ID:          strconv.FormatUint(c.ID, 10),
		UserID:      strconv.FormatUint(c.UserID, 10), // <- Key (not UserId)
		Name:        c.Name,
		Description: strPtr(c.Description),
		ColorHex:    strPtr(c.Color),
		Icon:        strPtr(c.Icon),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
