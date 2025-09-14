// Package controller provides mapping helpers between GraphQL models and core commands/domain for the Category context.
package controller

import (
	"strconv"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
)

// toModelOut converts a domain.Category to a GraphQL model.Category.
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

// toCreateCommand converts a GraphQL CreateCategoryInput into an input.CreateCategoryCommand.
func toCreateCommand(in gmodel.CreateCategoryInput, userID uint64) input.CreateCategoryCommand {
	return input.CreateCategoryCommand{
		Name:        in.Name,
		Description: in.Description,
		ColorHex:    in.ColorHex,
		Icon:        in.Icon,
		UserID:      userID,
	}
}

// toUpdateCommand converts a GraphQL UpdateCategoryInput into an input.UpdateCategoryCommand.
// It parses the ID field (string) into uint64 before passing it to the use case.
func toUpdateCommand(in gmodel.UpdateCategoryInput, userID uint64) (input.UpdateCategoryCommand, error) {
	id, err := strconv.ParseUint(in.ID, 10, 64)
	if err != nil {
		return input.UpdateCategoryCommand{}, err
	}
	return input.UpdateCategoryCommand{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		ColorHex:    in.ColorHex,
		Icon:        in.Icon,
		UserID:      userID,
	}, nil
}
