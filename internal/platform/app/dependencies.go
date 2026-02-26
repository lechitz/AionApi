// Package app defines the core application contracts.
package app

import (
	inputAdmin "github.com/lechitz/AionApi/internal/admin/core/ports/input"
	inputAudit "github.com/lechitz/AionApi/internal/audit/core/ports/input"
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
	AdminService    inputAdmin.AdminService
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	RecordService   inputRecord.RecordService
	ChatService     inputChat.ChatService
	AuditService    inputAudit.Service
	Logger          logger.ContextLogger
}
