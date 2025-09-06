package resolver

import (
	categoryController "github.com/lechitz/AionApi/internal/category/adapter/primary/graphql/handler"
	inputCategory "github.com/lechitz/AionApi/internal/category/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	tagsController "github.com/lechitz/AionApi/internal/tag/adapter/primary/graphql/handler"
	inputTag "github.com/lechitz/AionApi/internal/tag/core/ports/input"
)

// This file will not be regenerated automatically.
//
// It serves as a dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver.
type Resolver struct {
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	Logger          logger.ContextLogger

	controllers struct {
		category *categoryController.Handler
		tags     *tagsController.Handler
	}
}

// CategoryController returns a lazily constructed Category controller.
func (r *Resolver) CategoryController() *categoryController.Handler {
	if r.controllers.category == nil {
		r.controllers.category = categoryController.NewHandler(r.CategoryService, r.Logger)
	}
	return r.controllers.category
}

// TagsController returns a lazily built Tags handler (controller).
// Keep this for symmetry; if TagService isn't implemented yet,
// comment out usages in resolver until the service exists.
func (r *Resolver) TagsController() *tagsController.Handler {
	if r.controllers.tags == nil {
		r.controllers.tags = tagsController.NewHandler(r.TagService, r.Logger)
	}
	return r.controllers.tags
}
