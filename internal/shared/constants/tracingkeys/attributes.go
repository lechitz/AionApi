// Package tracingkeys defines shared keys for OpenTelemetry span attributes and audit fields.
package tracingkeys

// =============================================================================
// SEMANTIC ATTRIBUTES
// Based on OpenTelemetry Semantic Conventions
// =============================================================================

// Standard HTTP attributes (from OTEL spec).
const (
	AttrHTTPMethod     = "http.method"
	AttrHTTPStatusCode = "http.status_code"
	AttrHTTPRoute      = "http.route"
	AttrHTTPURL        = "http.url"
)

// Standard DB attributes (from OTEL spec).
const (
	AttrDBSystem    = "db.system"
	AttrDBOperation = "db.operation"
	AttrDBStatement = "db.statement"
)

// =============================================================================
// CUSTOM ATTRIBUTES - AionApi Domain
// Format: aion.<domain>.<attribute>
// =============================================================================

// Generic attributes.
const (
	AttrUserID    = "aion.user_id"
	AttrOperation = "aion.operation"
	AttrRequestIP = "aion.request.ip"
)

// Auth domain.
const (
	AttrAuthError  = "aion.auth.error"
	AttrAuthStatus = "aion.auth.status"
)

// Tag domain.
const (
	AttrTagID   = "aion.tag.id"
	AttrTagName = "aion.tag.name"
)

// Category domain.
const (
	AttrCategoryID   = "aion.category.id"
	AttrCategoryName = "aion.category.name"
)

// Record domain.
const (
	AttrRecordID    = "aion.record.id"
	AttrRecordTitle = "aion.record.title"
)

// Chat domain.
const (
	AttrChatMessageLength = "aion.chat.message_length"
	AttrChatTokensUsed    = "aion.chat.tokens_used"
)
