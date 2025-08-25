// Package constants defines constants used in authentication middleware.
package constants

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
	AttrAuthMiddlewareError  = "authmiddleware.error"
	AttrAuthMiddlewareUserID = "authmiddleware.userID"
	AttrAuthMiddlewareStatus = "authmiddleware.status"
)

// StatusAuthenticated is the status for authenticated span.
const StatusAuthenticated = "authenticated"

// MsgContextSet is the message for setting context in auth middleware.
const MsgContextSet = "authmiddleware context set"
