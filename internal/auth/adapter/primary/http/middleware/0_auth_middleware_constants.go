// Package middleware constants contains constants related to authentication operations.
package middleware

const (
	TracerAuthMiddleware    = "aionapi.auth.middleware"
	SpanAuthMiddleware      = "auth_middleware"
	SpanErrorMissingToken   = "missing_token"
	SpanErrorTokenInvalid   = "invalid_token"
	SpanStatusAuthenticated = "authenticated"

	AttrAuthMiddlewareError  = "auth_mw.error"
	AttrAuthMiddlewareUserID = "auth_mw.user_id"
	AttrAuthMiddlewareStatus = "auth_mw.status"

	StatusAuthenticated = "authenticated"

	ErrorUnauthorizedAccessMissingToken = "missing or empty auth token"
	ErrorUnauthorizedAccessInvalidToken = "invalid auth token"
	MsgContextSet                       = "auth context attached to request"
)
