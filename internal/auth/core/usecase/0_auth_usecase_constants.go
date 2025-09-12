// Package usecase constants contains constants related to authentication operations.
package usecase

// TracerName identifies the tracer used in auth use cases.
const TracerName = "aionapi.auth"

// Spans.
const (

	// SpanLogin is the name of the span for login.
	SpanLogin = "Login"

	// SpanLogout is the name of the span for logout.
	SpanLogout = "Logout"

	// SpanRevokeToken is the name of the span for revoke token.
	SpanRevokeToken = "RevokeToken"
)

// Errors.
const (

	// ErrorToCompareHashAndPassword is the error message when the password is invalid.
	ErrorToCompareHashAndPassword = "invalid credentials" // #nosec G101

	// ErrorToCreateToken is the error message when the token cannot be created. // #nosec G101.
	ErrorToCreateToken = "error to create token"

	// ErrorToRevokeToken is the error message when the token cannot be revoked. // #nosec G101.
	ErrorToRevokeToken = "error to revoke token"

	// ErrorToDeleteToken is the error message when the token cannot be deleted. // #nosec G101.
	ErrorToDeleteToken = "error to delete token"

	// ErrorToGetUserByUserName is the error message when the user cannot be found. // #nosec G101.
	ErrorToGetUserByUserName = "error to get user by username"

	// UserNotFoundOrInvalidCredentials is the error message when the user cannot be found.
	UserNotFoundOrInvalidCredentials = "user not found or invalid credentials"

	// InvalidCredentials is the error message when the credentials are invalid.
	InvalidCredentials = "invalid credentials" // #nosec G101
)

// Success messages.
const (
	// SuccessToLogin is the success message when the user logs in.
	SuccessToLogin = "user logged in successfully"

	// SuccessUserLoggedOut is the success message when the user logs out.
	SuccessUserLoggedOut = "user logged out successfully"
)

// Events (trace).
const (

	// EventLookupUser is emitted right before looking up the user.
	EventLookupUser = "lookup_user"

	// EventComparePassword is emitted right before comparing the password.
	EventComparePassword = "compare_password"

	// EventGenerateToken is emitted right before generating the token.
	EventGenerateToken = "generate_token"

	// EventSaveTokenToStore is emitted right before saving the token to the cache.
	EventSaveTokenToStore = "save_token_to_store"

	// EventTokenRevoked is emitted right before revoking the token.
	EventTokenRevoked = "token_revoked"

	// EventLoginSuccess is emitted right after successful login.
	EventLoginSuccess = "login_success"

	// EventRevokeToken is emitted right before revoking the token.
	EventRevokeToken = "revoke_token"

	// EventLogoutSuccess is emitted right after a successful logout.
	EventLogoutSuccess = "logout_success"
)
