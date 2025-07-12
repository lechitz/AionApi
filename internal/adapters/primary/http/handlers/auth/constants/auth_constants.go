// Package constants contains constants used throughout the auth handler.
package constants

// TracerAuthHandler is the tracer name for auth handler operations in OpenTelemetry.
const TracerAuthHandler = "aionapi.auth.handler"

// Span names for OpenTelemetry auth handler operations.
const (
	SpanLoginHandler  = "auth.login"
	SpanLogoutHandler = "auth.logout"
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
	ErrMissingToken  = "missing token in context"
	ErrUnauthorized  = "unauthorized"
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
	LogMissingToken  = "missing token"
	LogLogoutFailed  = "logout failed"
)
