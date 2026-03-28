// Package middleware constants contains constants related to authentication operations.
package middleware

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerAuthMiddleware is the tracer name used by the auth middleware.
// Format: aion-api.<domain>.<layer> .
const TracerAuthMiddleware = "aion-api.auth.middleware"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

const (
	// SpanAuthMiddleware is the span name for the auth middleware.
	SpanAuthMiddleware = "auth.middleware"

	// SpanErrorMissingToken is the span name for missing token errors.
	SpanErrorMissingToken = "auth.middleware.missing_token" // #nosec G101: span name, not a credential

	// SpanErrorTokenInvalid is the span name for invalid token errors.
	SpanErrorTokenInvalid = "auth.middleware.invalid_token" // #nosec G101: span name, not a credential

	// SpanStatusAuthenticated is the span status for successful authentication.
	SpanStatusAuthenticated = "authenticated"
)

// -----------------------------------------------------------------------------
// Span Attributes
// Format: aion.<domain>.<attribute>
// -----------------------------------------------------------------------------

const (
	// AttrAuthMiddlewareError is the attribute key used to log auth middleware errors.
	AttrAuthMiddlewareError = "aion.auth.error"

	// AttrAuthMiddlewareUserID is the attribute key used to log auth middleware user ID.
	AttrAuthMiddlewareUserID = "aion.user_id"

	// AttrAuthMiddlewareStatus is the attribute key used to log auth middleware status.
	AttrAuthMiddlewareStatus = "aion.auth.status"
)

// =============================================================================
// BUSINESS LOGIC - Status and Messages
// =============================================================================

const (
	// StatusAuthenticated is the status string used when a request is authenticated.
	StatusAuthenticated = "authenticated"

	// MsgS2SAuthBypass is the log message used when S2S authentication bypass occurs.
	MsgS2SAuthBypass = "S2S auth bypass: service account authenticated" // #nosec G101: false positive — log message, not credentials

	// ErrorUnauthorizedAccessMissingToken is the user-facing error when auth token is missing.
	ErrorUnauthorizedAccessMissingToken = "missing or empty auth token" // #nosec G101: false positive — user-facing error message, not a credential

	// ErrorUnauthorizedAccessInvalidToken is the user-facing error when auth token is invalid.
	ErrorUnauthorizedAccessInvalidToken = "invalid auth token" // #nosec G101: false positive — user-facing error message, not a credential

	// MsgContextSet is the log message used when auth context is attached to the request.
	MsgContextSet = "auth context attached to request"
)
