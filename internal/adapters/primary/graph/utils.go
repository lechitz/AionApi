// Package graph The package graph Contains GraphQL resolvers and helpers for the API.
package graph

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/primary/graph/model"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/ctxkeys"
)

// toGraphQLCategory converts a domain.Category to a GraphQL model.Category.
func toGraphQLCategory(c domain.Category) *model.Category {
	return &model.Category{
		CategoryID:  strconv.FormatUint(c.ID, 10),
		UserID:      strconv.FormatUint(c.UserID, 10),
		Name:        c.Name,
		Description: &c.Description,
		ColorHex:    &c.Color,
		Icon:        &c.Icon,
	}
}

// getClientMeta extracts client IP and User-Agent from context.
func getClientMeta(ctx context.Context) (string, string) {
	ip, _ := ctx.Value(ctxkeys.RequestIP).(string)
	userAgent, _ := ctx.Value(ctxkeys.RequestUserAgent).(string)
	return ip, userAgent
}
