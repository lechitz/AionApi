package resolvers

import (
	categoryController "github.com/lechitz/AionApi/internal/adapters/primary/graph/controllers/category"
	tagsController "github.com/lechitz/AionApi/internal/adapters/primary/graph/controllers/tags"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// This file will not be regenerated automatically.
//
// It serves as a dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver.
type Resolver struct {
	CategoryService input.CategoryService
	TagService      input.TagService
	Logger          output.ContextLogger

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
// comment out usages in resolvers until the service exists.
func (r *Resolver) TagsController() *tagsController.Handler {
	if r.controllers.tags == nil {
		r.controllers.tags = tagsController.NewHandler(r.TagService, r.Logger)
	}
	return r.controllers.tags
}
