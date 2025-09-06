// Package middleware constants contains constants related to authentication operations.
package middleware

// Tracers
const (
	// TracerAuthMiddleware is the tracer name for auth middleware.
	TracerAuthMiddleware = "AuthMiddleware"

	// SpanAuthMiddleware is the span name for auth middleware.
	SpanAuthMiddleware = "Auth"
)

// Span status messages
const (
	SpanErrorMissingToken   = "missing or invalid token"
	SpanStatusAuthenticated = "authenticated"
	SpanErrorTokenInvalid   = "token invalid"
)

// Attribute keys
const (
	AttrAuthMiddlewareError  = "middleware.error"
	AttrAuthMiddlewareUserID = "middleware.userID"
	AttrAuthMiddlewareStatus = "middleware.status"
)

// StatusAuthenticated is the status for authenticated span.
const StatusAuthenticated = "authenticated"

// MsgContextSet is the message for setting context in auth middleware.
const MsgContextSet = "middleware context set"

// ErrorUnauthorizedAccessMissingToken is returned when no authentication token is present in the request.
const ErrorUnauthorizedAccessMissingToken = "unauthorized access: missing token"

// ErrorUnauthorizedAccessInvalidToken is returned when the authentication token provided is invalid.
const ErrorUnauthorizedAccessInvalidToken = "unauthorized access: invalid token"
