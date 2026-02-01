// Package tracingkeys defines shared keys for OpenTelemetry span attributes and audit fields.
package tracingkeys

// =============================================================================
// SPAN STATUS DESCRIPTIONS
// Used with span.SetStatus()
// =============================================================================

const (
	// StatusCreated indicates a resource was created successfully.
	StatusCreated = "created"
	// StatusRetrieved indicates a resource was retrieved successfully.
	StatusRetrieved = "retrieved"
	// StatusUpdated indicates a resource was updated successfully.
	StatusUpdated = "updated"
	// StatusDeleted indicates a resource was deleted successfully.
	StatusDeleted = "deleted"
	// StatusAuthenticated indicates successful authentication.
	StatusAuthenticated = "authenticated"
	// StatusValidated indicates successful validation.
	StatusValidated = "validated"
)

// =============================================================================
// HTTP/REQUEST ATTRIBUTES (Legacy - use attributes.go for new code)
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
