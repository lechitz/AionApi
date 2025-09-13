package graphql

import (
	categoryController "github.com/lechitz/AionApi/internal/category/adapter/primary/graphql/controller"
	inputCategory "github.com/lechitz/AionApi/internal/category/core/ports/input"
	inputTag "github.com/lechitz/AionApi/internal/tag/core/ports/input"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Resolver wires services into thin GraphQL controllers per context.
type Resolver struct {
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	Logger          logger.ContextLogger
}

// CategoryController returns the GraphQL controller for Category.
func (r *Resolver) CategoryController() categoryController.controller {
	return categoryController.NewController(r.CategoryService, r.Logger)
}

// TODO: Nota: Quando o contexto de Tags tiver um controller similar ao de Category, preciso adiciona-lo aqui:
//// TagController returns the GraphQL controller for Tag.
// func (r *Resolver) TagController() tagController.controller {
//	return tagController.NewController(r.TagService, r.Logger)
// }
