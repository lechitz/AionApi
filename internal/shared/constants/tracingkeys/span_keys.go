// Package tracingkeys defines shared keys for OpenTelemetry span attributes.
//
// IMPORTANT: This package contains LEGACY constants for HTTP/Request tracing only.
//
// Per AGENTS.md: "No magic strings: shared constants in internal/shared/constants;
// local constants in *_constants.go within the package."
//
// Each bounded context should define its own tracing constants locally in
// internal/<ctx>/core/usecase/*_constants.go including:
//   - TracerName (e.g., "aionapi.category.usecase")
//   - Span names (e.g., "category.create")
//   - Event names (e.g., "category.cache.hit")
//   - Status descriptions (e.g., "created", "retrieved")
//
// Examples:
//   - internal/category/core/usecase/0_category_usecase_constants.go
//   - internal/tag/core/usecase/0_tag_usecase_constants.go
//   - internal/record/core/usecase/0_record_usecase_constants.go
package tracingkeys

// =============================================================================
// HTTP/REQUEST ATTRIBUTES
// Legacy constants kept for backward compatibility with existing HTTP handlers.
// =============================================================================

const (
	// HTTPStatusCodeKey is the key for setting the HTTP status code in tracing attributes.
	HTTPStatusCodeKey = "http.status_code"

	// RequestIPKey is the key for recording the request's originating IP address.
	RequestIPKey = "request.ip"

	// RequestUserAgentKey is the key for recording the request's User-Agent header.
	RequestUserAgentKey = "request.user_agent"

	// ContextUserID is the key for recording the user ID from the context.
	ContextUserID = "context.user_id"
)

// =============================================================================
// ERROR/DEBUG ATTRIBUTES
// =============================================================================

const (
	// RecoveredKey is used for logging recovered panics.
	RecoveredKey = "recovered"

	// StacktraceKey is the key for recording the stack trace of an error.
	StacktraceKey = "stacktrace"

	// ErrorID is the key for recording the error ID.
	ErrorID = "error.id"
)
