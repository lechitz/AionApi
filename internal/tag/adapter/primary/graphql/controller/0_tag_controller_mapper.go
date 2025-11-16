// Package controller provides mapping helpers between GraphQL models and core commands/domain for the Category context.
package controller

import (
	"strconv"
	"time"

	gmodel "github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

// toModelOut converts a domain.Tag to a GraphQL model.Tag.
func toModelOut(t domain.Tag) *gmodel.Tag {
	out := &gmodel.Tag{
		ID:         strconv.FormatUint(t.ID, 10),
		UserID:     strconv.FormatUint(t.UserID, 10),
		CategoryID: strconv.FormatUint(t.CategoryID, 10),
		Name:       t.Name,
	}
	if t.Description != "" {
		out.Description = &t.Description
	}

	out.CreatedAt = t.CreatedAt.Format(time.RFC3339)
	out.UpdatedAt = t.UpdatedAt.Format(time.RFC3339)

	return out
}

// toCreateCommand converts a GraphQL CreateTagInput into an input.CreateTagCommand.
func toCreateCommand(in gmodel.CreateTagInput, userID, categoryID uint64) input.CreateTagCommand {
	return input.CreateTagCommand{
		Name:        in.Name,
		Description: in.Description,
		UserID:      userID,
		CategoryID:  categoryID,
	}
}
