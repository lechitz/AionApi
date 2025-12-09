// Package tracingkeys defines shared keys for OpenTelemetry span attributes and audit fields.
package tracingkeys

// =============================================================================
// TRACER NAMES - Instrumentation Scopes
// Format: aionapi.<domain>.<layer>
// =============================================================================

// Auth domain.
const (
	TracerAuthUsecase    = "aionapi.auth.usecase"
	TracerAuthMiddleware = "aionapi.auth.middleware"
	TracerAuthHandler    = "aionapi.auth.handler"
	TracerAuthCache      = "aionapi.auth.cache"
)

// User domain.
const (
	TracerUserUsecase    = "aionapi.user.usecase"
	TracerUserHandler    = "aionapi.user.handler"
	TracerUserRepository = "aionapi.user.repository"
)

// Tag domain.
const (
	TracerTagController = "aionapi.tag.controller"
	TracerTagUsecase    = "aionapi.tag.usecase"
	TracerTagRepository = "aionapi.tag.repository"
)

// Category domain.
const (
	TracerCategoryController = "aionapi.category.controller"
	TracerCategoryUsecase    = "aionapi.category.usecase"
	TracerCategoryRepository = "aionapi.category.repository"
)

// Record domain.
const (
	TracerRecordController = "aionapi.record.controller"
	TracerRecordUsecase    = "aionapi.record.usecase"
	TracerRecordRepository = "aionapi.record.repository"
)

// Chat domain.
const (
	TracerChatHandler = "aionapi.chat.handler"
	TracerChatUsecase = "aionapi.chat.usecase"
	TracerChatClient  = "aionapi.chat.client"
)

// Admin domain.
const (
	TracerAdminHandler = "aionapi.admin.handler"
	TracerAdminUsecase = "aionapi.admin.usecase"
)
