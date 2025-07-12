// Package ctxkeys contains constants used for context keys.
package ctxkeys

// contextKey is a custom string type for context keys.
type contextKey string

// Context keys for context values.
const (
	UserID           contextKey = "user_id"            // Context key for user ID.
	Token            contextKey = "token"              // Context key for auth token.
	RequestID        contextKey = "request_id"         // Context key for request ID.
	TraceID          contextKey = "trace_id"           // Context key for trace ID.
	SpanID           contextKey = "span_id"            // Context key for span ID.
	RequestIP        contextKey = "request_ip"         // Context key for request IP.
	RequestUserAgent contextKey = "request_user_agent" // Context key for request user agent.
)
