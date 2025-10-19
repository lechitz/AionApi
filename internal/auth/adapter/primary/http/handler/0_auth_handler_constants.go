// Package handler contains constants used throughout the auth handler.middleware
package handler

// TracerAuthHandler is the tracer name for auth handler operations in OpenTelemetry.
const TracerAuthHandler = "aionapi.auth.handler"

// Span names for OpenTelemetry auth handler operations.
const (
	SpanLoginHandler   = "auth.login"
	SpanLogoutHandler  = "auth.logout"
	SpanRefreshHandler = "auth.refresh"
)

// Attribute keys used in spans for auth handlers.
const (
	AttrRefreshTokenPresent = "refresh_token_present"
)

// Event names for key points within auth handler spans.
const (
	EventDecodeRequest     = "decode_request"
	EventAuthServiceLogin  = "auth_service.login"
	EventAuthServiceLogout = "auth_service.logout"
	EventLoginSuccess      = "auth.login.success"
	EventLogoutSuccess     = "auth.logout.success"
)

// Status names for semantic span states.
const (
	StatusLoginSuccess  = "login_success"
	StatusLogoutSuccess = "logout_success"
)

// Error messages used in auth handler (for response, tracing and logs).
const (
	ErrMissingUserID = "missing user id in context"
	ErrLogin         = "error on login"
	ErrLogout        = "error on logout"
)

// Success messages used in auth handler.
const (
	MsgLoginSuccess  = "user logged in successfully"
	MsgLogoutSuccess = "user logged out successfully"
)

// Log message keys for structured and leveled logging.
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
