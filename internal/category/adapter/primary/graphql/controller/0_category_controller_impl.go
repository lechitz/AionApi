package controller

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/primary/graphql/model"
	"github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// CategoryController is the contract used by GraphQL resolvers.
// Keep it thin: map GraphQL <-> domain, add tracing/logging, delegate to use cases.
type CategoryController interface {
	Create(ctx context.Context, in model.CreateCategoryInput, userID uint64) (*model.Category, error)
	Update(ctx context.Context, in model.UpdateCategoryInput, userID uint64) (*model.Category, error)
	SoftDelete(ctx context.Context, categoryID, userID uint64) error

	GetByID(ctx context.Context, categoryID uint64, userID uint64) (*model.Category, error)
	GetByName(ctx context.Context, categoryName string, userID uint64) (*model.Category, error)
	ListAll(ctx context.Context, userID uint64) ([]*model.Category, error)
}

// controller is the concrete implementation of CategoryController.
type controller struct {
	CategoryService input.CategoryService
	Logger          logger.ContextLogger
}

// NewController wires dependencies and returns a CategoryController.
func NewController(svc input.CategoryService, log logger.ContextLogger) CategoryController {
	return &controller{
		CategoryService: svc,
		Logger:          log,
	}
}
