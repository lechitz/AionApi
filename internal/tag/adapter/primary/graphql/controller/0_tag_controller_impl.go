package controller

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

// TagController is the contract used by GraphQL resolvers.
// Keep it thin: map GraphQL <-> domain, add tracing/logging, delegate to use cases.
type TagController interface {
	Create(ctx context.Context, in model.CreateTagInput, userID uint64) (*model.Tag, error)
	GetByID(ctx context.Context, tagID, userID uint64) (*model.Tag, error)
	GetByName(ctx context.Context, tagName string, userID uint64) (*model.Tag, error)
}

// controller is the controller for the tag service.
type controller struct {
	TagService input.TagService
	Logger     logger.ContextLogger
}

// NewController wires dependencies and returns a Controller.
func NewController(svc input.TagService, logger logger.ContextLogger) TagController {
	return &controller{
		TagService: svc,
		Logger:     logger,
	}
}
