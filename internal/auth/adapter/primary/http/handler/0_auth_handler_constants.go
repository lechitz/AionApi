// Package handler contains constants used throughout the auth handler.
package handler

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerAuthHandler is the tracer name for auth handler operations in OpenTelemetry.
// Format: aionapi.<domain>.<layer> .
const TracerAuthHandler = "aionapi.auth.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

// Span names for auth handler operations.
const (
	SpanLoginHandler   = "auth.handler.login"
	SpanLogoutHandler  = "auth.handler.logout"
	SpanRefreshHandler = "auth.handler.refresh"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

// Event names for auth handler tracing.
const (
	EventDecodeRequest     = "auth.handler.decode_request"
	EventAuthServiceLogin  = "auth.handler.service_login"
	EventAuthServiceLogout = "auth.handler.service_logout"
	EventLoginSuccess      = "auth.handler.login_success"
	EventLogoutSuccess     = "auth.handler.logout_success"
)

// -----------------------------------------------------------------------------
// Status Names
// -----------------------------------------------------------------------------

// Status descriptions for auth handler spans.
const (
	StatusLoginSuccess  = "login_success"
	StatusLogoutSuccess = "logout_success"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	ErrMissingUserID = "missing user id in context"
	ErrLogin         = "error on login"
	ErrLogout        = "error on logout"
)

// Success messages.
const (
	MsgLoginSuccess  = "user logged in successfully"
	MsgLogoutSuccess = "user logged out successfully"
)

// Log message keys.
const (
	LogMissingUserID = "missing user id"
	LogLogoutFailed  = "logout failed"
)

// Validation errors.
const (
	ErrRequiredFields    = "username and password are required"
	ErrMinUsernameLength = "username must have at least 3 characters"
	ErrMinPasswordLength = "password must have at least 8 characters"
	MinUsernameLength    = 3
	MinPasswordLength    = 8
)
