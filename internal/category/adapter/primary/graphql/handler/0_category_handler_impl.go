package handler

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Controller is the contract used by GraphQL resolvers.
// Keep it thin: map GraphQL <-> domain, add tracing/logging, delegate to use cases.
type Controller interface {
	Create(ctx context.Context, in model.CreateCategoryInput, userID uint64) (*model.Category, error)
	Update(ctx context.Context, in model.UpdateCategoryInput, userID uint64) (*model.Category, error)
	SoftDelete(ctx context.Context, categoryID, userID uint64) error

	GetByID(ctx context.Context, categoryID uint64, userID uint64) (*model.Category, error)
	GetByName(ctx context.Context, categoryName string, userID uint64) (*model.Category, error)
	ListAll(ctx context.Context, userID uint64) ([]*model.Category, error)
}

// Handler is the concrete implementation of Controller.
type Handler struct {
	CategoryService input.CategoryService
	Logger          logger.ContextLogger
}

// NewHandler wires dependencies and returns a Controller.
func NewHandler(svc input.CategoryService, log logger.ContextLogger) Controller {
	return &Handler{
		CategoryService: svc,
		Logger:          log,
	}
}
