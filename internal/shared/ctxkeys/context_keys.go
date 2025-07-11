// Package ctxkeys contains constants used for context keys.
package ctxkeys

// contextKey is a string type used for context keys.
type contextKey string

// UserID is a constant string representing the key used to define or identify a user ID.
const (
	UserID contextKey = "user_id"

	// Token is a constant string representing the key used to define or identify a token.
	Token contextKey = "token"

	CtxKeyRequestID contextKey = "request_id"
	CtxKeyTraceID   contextKey = "trace_id"
	CtxKeyUserID    contextKey = "user_id"
)
