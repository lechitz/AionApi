// Package app defines the core application contracts.
package app

import (
	inputAuth "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	inputCategory "github.com/lechitz/AionApi/internal/category/core/ports/input"
	inputChat "github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	inputRecord "github.com/lechitz/AionApi/internal/record/core/ports/input"
	inputTag "github.com/lechitz/AionApi/internal/tag/core/ports/input"
	inputUser "github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// Dependencies exposes application services that primary adapters (HTTP/GraphQL) consume.
// This is the contract between the application layer and presentation layer.
type Dependencies struct {
	AuthService     inputAuth.AuthService
	UserService     inputUser.UserService
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	RecordService   inputRecord.RecordService
	ChatService     inputChat.ChatService
	Logger          logger.ContextLogger
}
