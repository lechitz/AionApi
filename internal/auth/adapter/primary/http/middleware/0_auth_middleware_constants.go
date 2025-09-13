// Package middleware constants contains constants related to authentication operations.
package middleware

const (
	// TracerAuthMiddleware is the name of the tracer used in auth middleware.
	TracerAuthMiddleware = "aionapi.auth.middleware"

	// SpanAuthMiddleware is the name of the span for auth middleware.
	SpanAuthMiddleware = "auth_middleware"

	// SpanErrorMissingToken is the name of the span for the missing token.
	SpanErrorMissingToken = "missing_token"

	// SpanErrorTokenInvalid is the name of the span for an invalid token.
	SpanErrorTokenInvalid = "invalid_token"

	// SpanStatusAuthenticated is the name of the span for authenticated.
	SpanStatusAuthenticated = "authenticated"

	// AttrAuthMiddlewareError is the name of the attribute for auth middleware error.
	AttrAuthMiddlewareError = "auth_mw.error"

	// AttrAuthMiddlewareUserID is the name of the attribute for auth middleware user ID.
	AttrAuthMiddlewareUserID = "auth_mw.user_id"

	// AttrAuthMiddlewareStatus is the name of the attribute for auth middleware status.
	AttrAuthMiddlewareStatus = "auth_mw.status"

	// StatusAuthenticated is the status for authenticated.
	StatusAuthenticated = "authenticated"

	// ErrorUnauthorizedAccessMissingToken is the error message for a missing or empty auth token.
	ErrorUnauthorizedAccessMissingToken = "missing or empty auth token" // #nosec G101: false positive — user-facing error message, not a credential

	// ErrorUnauthorizedAccessInvalidToken is the error message for an invalid auth token.
	ErrorUnauthorizedAccessInvalidToken = "invalid auth token" // #nosec G101: false positive — user-facing error message, not a credential

	// MsgContextSet is the message for when the auth context is set.
	MsgContextSet = "auth context attached to request"
)
