// Package ctxkeys contains constants used for context keys.
package ctxkeys

// contextKey is a custom string type for context keys.
type contextKey string

// Context keys for context values.
const (
	UserID         contextKey = "user_id"         // Context key for user ID.
	Token          contextKey = "token"           // Context key for auth token.
	RequestID      contextKey = "request_id"      // Context key for request ID.
	TraceID        contextKey = "trace_id"        // Context key for trace ID.
	SpanID         contextKey = "span_id"         // Context key for span ID.
	Claims         contextKey = "claims"          // Context key for claims.
	IdempotencyKey contextKey = "idempotency_key" // Context key for Idempotency-Key header.
)
