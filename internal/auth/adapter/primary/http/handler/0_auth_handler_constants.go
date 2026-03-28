// Package handler contains constants used throughout the auth handler.
package handler

// =============================================================================
// TRACING - OpenTelemetry Instrumentation
// =============================================================================

// TracerAuthHandler is the tracer name for auth handler operations in OpenTelemetry.
// Format: aion-api.<domain>.<layer> .
const TracerAuthHandler = "aion-api.auth.handler"

// -----------------------------------------------------------------------------
// Span Names
// Format: <domain>.<operation>
// -----------------------------------------------------------------------------

// Span names for auth handler operations.
const (
	SpanLoginHandler   = "auth.handler.login"
	SpanLogoutHandler  = "auth.handler.logout"
	SpanRefreshHandler = "auth.handler.refresh"
	SpanSessionHandler = "auth.handler.session"
)

// -----------------------------------------------------------------------------
// Event Names
// Format: <domain>.<action>.<detail>
// -----------------------------------------------------------------------------

// Event names for auth handler tracing.
const (
	EventDecodeRequest      = "auth.handler.decode_request"
	EventAuthServiceLogin   = "auth.handler.service_login"
	EventAuthServiceLogout  = "auth.handler.service_logout"
	EventAuthServiceRefresh = "auth.handler.service_refresh"
	EventLoginSuccess       = "auth.handler.login_success"
	EventLogoutSuccess      = "auth.handler.logout_success"
	EventRefreshSuccess     = "auth.handler.refresh_success"
	EventSessionSuccess     = "auth.handler.session_success"
)

// -----------------------------------------------------------------------------
// Status Names
// -----------------------------------------------------------------------------

// Status descriptions for auth handler spans.
const (
	StatusLoginSuccess   = "login_success"
	StatusLogoutSuccess  = "logout_success"
	StatusRefreshSuccess = "refresh_success"
	StatusSessionSuccess = "session_success"
)

// =============================================================================
// BUSINESS LOGIC - Error and Success Messages
// =============================================================================

// Error messages.
const (
	ErrMissingUserID = "missing user id in context"
	ErrLogin         = "error on login"
	ErrLogout        = "error on logout"
	ErrRefresh       = "error on refresh"
	ErrSession       = "error on session"
)

// Success messages.
const (
	MsgLoginSuccess   = "user logged in successfully"
	MsgLogoutSuccess  = "user logged out successfully"
	MsgRefreshSuccess = "token refreshed successfully"
	MsgSessionSuccess = "session fetched successfully"
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

// =============================================================================
// HTTP - Request limits and validation field names
// =============================================================================

const (
	// LoginMaxBodyBytes limits the /auth/login payload size.
	LoginMaxBodyBytes int64 = 1 << 20 // 1MB

	// ValidationFieldCredentials is the field name used in validation errors for login credentials.
	ValidationFieldCredentials = "credentials"
)

// =============================================================================
// SESSION - Errors and tracing attribute keys
// =============================================================================

const (
	// ErrMissingAccessToken is returned when access cookie is missing or empty.
	ErrMissingAccessToken = "missing access token"

	// AttrAccessTokenPresent records whether access token cookie was present. Never store the token value.
	AttrAccessTokenPresent = "access_token_present" // #nosec G101: attribute name, not a credential
)

// =============================================================================
// REFRESH - Errors and tracing attribute keys
// =============================================================================

const (
	// ErrMissingRefreshToken is returned when refresh cookie is missing or empty.
	ErrMissingRefreshToken = "missing refresh token"

	// AttrRefreshTokenPresent records whether refresh token cookie was present. Never store the token value.
	AttrRefreshTokenPresent = "refresh_token_present" // #nosec G101: attribute name, not a credential
)
