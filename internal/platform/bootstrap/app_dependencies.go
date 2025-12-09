// Package bootstrap holds shared domain dependency contracts used by adapters.
package bootstrap

import (
	inputAuth "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	inputCategory "github.com/lechitz/AionApi/internal/category/core/ports/input"
	inputChat "github.com/lechitz/AionApi/internal/chat/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	inputRecord "github.com/lechitz/AionApi/internal/record/core/ports/input"
	inputTag "github.com/lechitz/AionApi/internal/tag/core/ports/input"
	inputUser "github.com/lechitz/AionApi/internal/user/core/ports/input"
)

// AppDependencies exposes input ports that primary adapters consume.
// Retained for HTTP/GraphQL composition while the Fx bootstrap handles wiring.
type AppDependencies struct {
	AuthService     inputAuth.AuthService
	UserService     inputUser.UserService
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	RecordService   inputRecord.RecordService
	ChatService     inputChat.ChatService

	Logger logger.ContextLogger
}
