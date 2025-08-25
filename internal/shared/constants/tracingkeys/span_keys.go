// Package tracingkeys defines shared keys for OpenTelemetry span attributes and audit fields.
package tracingkeys

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

const (
	// RecoveredKey and StacktraceKey are used for logging errors.
	RecoveredKey = "recovered"

	// StacktraceKey is the key for recording the stack trace of an error.
	StacktraceKey = "stacktrace"
)

// ErrorID is the key for recording the error ID.
const ErrorID = "error.id"
