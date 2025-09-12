// Package handler contains the GraphQL handlers for the category service.
package handler

import (
	"strconv"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/core/domain"
)

// toModelOut is a helper function to convert a domain.Category to a gmodel.Category.
func toModelOut(c domain.Category) *gmodel.Category {
	out := &gmodel.Category{
		ID:     strconv.FormatUint(c.ID, 10),
		UserID: strconv.FormatUint(c.UserID, 10),
		Name:   c.Name,
	}
	if c.Description != "" {
		out.Description = &c.Description
	}
	if c.Color != "" {
		out.ColorHex = &c.Color
	}
	if c.Icon != "" {
		out.Icon = &c.Icon
	}
	return out
}

// toDomainCreate is a helper function to convert a gmodel.CreateCategoryInput to a domain.Category.
func toDomainCreate(in gmodel.CreateCategoryInput, userID uint64) domain.Category {
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

// toDomainUpdate is a helper function to convert a gmodel.UpdateCategoryInput to a domain.Category.
func toDomainUpdate(in gmodel.UpdateCategoryInput, userID uint64) (domain.Category, error) {
	id, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return domain.Category{}, err
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
	return c, nil
}
