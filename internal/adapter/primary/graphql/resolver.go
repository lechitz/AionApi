package graphql

import (
	categoryController "github.com/lechitz/aion-api/internal/category/adapter/primary/graphql/controller"
	chatController "github.com/lechitz/aion-api/internal/chat/adapter/primary/graphql/controller"
	recordController "github.com/lechitz/aion-api/internal/record/adapter/primary/graphql/controller"
	tagController "github.com/lechitz/aion-api/internal/tag/adapter/primary/graphql/controller"

	categoryInput "github.com/lechitz/aion-api/internal/category/core/ports/input"
	chatInput "github.com/lechitz/aion-api/internal/chat/core/ports/input"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	recordInput "github.com/lechitz/aion-api/internal/record/core/ports/input"
	tagInput "github.com/lechitz/aion-api/internal/tag/core/ports/input"
	userInput "github.com/lechitz/aion-api/internal/user/core/ports/input"
)

// Resolver wires services into thin GraphQL controllers per context.
type Resolver struct {
	CategoryService categoryInput.CategoryService
	TagService      tagInput.TagService
	RecordService   recordInput.RecordService
	ChatService     chatInput.ChatService
	UserService     userInput.UserService
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

// ChatController returns the Chat adapter controller (interface).
func (r *Resolver) ChatController() chatController.ChatController {
	return chatController.NewController(r.ChatService, r.Logger)
}
