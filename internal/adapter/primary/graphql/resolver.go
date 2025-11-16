package graphql

import (
	categoryController "github.com/lechitz/AionApi/internal/category/adapter/primary/graphql/controller"
	recordController "github.com/lechitz/AionApi/internal/record/adapter/primary/graphql/controller"
	tagController "github.com/lechitz/AionApi/internal/tag/adapter/primary/graphql/controller"

	categoryInput "github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	recordInput "github.com/lechitz/AionApi/internal/record/core/ports/input"
	tagInput "github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

// Resolver wires services into thin GraphQL controllers per context.
type Resolver struct {
	CategoryService categoryInput.CategoryService
	TagService      tagInput.TagService
	RecordService   recordInput.RecordService
	Logger          logger.ContextLogger
}

// CategoryController returns the Category adapter controller (interface).
func (r *Resolver) CategoryController() categoryController.CategoryController {
	return categoryController.NewController(r.CategoryService, r.Logger)
}

// TagController returns  the Tag adapter controller (interface).
func (r *Resolver) TagController() tagController.TagController {
	return tagController.NewController(r.TagService, r.Logger)
}

// RecordController returns the Record adapter controller (interface).
func (r *Resolver) RecordController() recordController.RecordController {
	return recordController.NewController(r.RecordService, r.Logger)
}
