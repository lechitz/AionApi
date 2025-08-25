// Package constants contains constants related to authentication operations.
package constants

// TracerName identifies the tracer used in auth use cases.
const TracerName = "aionapi.auth"

// Spans
const (
	SpanLogin  = "Login"
	SpanLogout = "Logout"
)

// Errors
const (
	ErrorToCompareHashAndPassword = "invalid credentials" // #nosec G101

	ErrorToCreateToken       = "error to create token"
	ErrorToRevokeToken       = "error to revoke token"
	ErrorToGetUserByUserName = "error to get user by username"

	UserNotFoundOrInvalidCredentials = "user not found or invalid credentials"
	InvalidCredentials               = "invalid credentials" // #nosec G101
)

// Success messages
const (
	SuccessToLogin       = "user logged in successfully"
	SuccessUserLoggedOut = "user logged out successfully"
)

// Events (trace)
const (
	EventLookupUser       = "lookup_user"
	EventComparePassword  = "compare_password"
	EventGenerateToken    = "generate_token"
	EventSaveTokenToStore = "save_token_to_store"
	EventLoginSuccess     = "login_success"

	EventRevokeToken   = "revoke_token"
	EventLogoutSuccess = "logout_success"
)
