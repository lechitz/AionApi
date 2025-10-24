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

// SpanValidateToken is the name of the span for validating a token.
const SpanValidateToken = "ValidateToken"

// ErrorInvalidToken indicates the token is invalid.
const ErrorInvalidToken = "invalid access reference"

// ErrorInvalidUserIDClaim indicates the user ID in claimsextractor is invalid.
const ErrorInvalidUserIDClaim = "invalid userID in claimsextractor"

// ErrorToRetrieveTokenFromCache indicates an error retrieving the access reference from cache.
const ErrorToRetrieveTokenFromCache = "error retrieving access reference from cache"

// ErrorTokenMismatch indicates the provided token does not match the stored one.
const ErrorTokenMismatch = "provided reference does not match stored one"

// SuccessTokenValidated indicates the access reference was validated successfully.
const SuccessTokenValidated = "access reference validated successfully"

// Errors.
const (

	// ErrorToCompareHashAndPassword is the error message when the password is invalid.
	ErrorToCompareHashAndPassword = "invalid credentials" // #nosec G101

	// ErrorToCreateToken is the error message when the token cannot be created. // #nosec G101.
	ErrorToCreateToken = "error to create token"

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

	// EventComparePassword //nolint:gosec,goconst is emitted right before comparing the password.
	EventComparePassword = "compare_password"

	// EventGenerateToken is emitted right before generating the token.
	EventGenerateToken = "generate_token"

	// EventSaveTokenToStore is emitted right before saving the token to the cache.
	EventSaveTokenToStore = "save_token_to_store" // #nosec G101: false positive — event name string, not a credential

	// EventTokenRevoked is emitted right before revoking the token.
	EventTokenRevoked = "token_revoked"

	// EventLoginSuccess is emitted right after successful login.
	EventLoginSuccess = "login_success"

	// EventRevokeToken is emitted right before revoking the token.
	EventRevokeToken = "revoke_token"

	// EventVerifyToken is emitted right before verifying token signature/exp.
	EventVerifyToken = "verify_token"

	// EventExtractUserID is emitted before extracting/parsing the userID claim.
	EventExtractUserID = "extract_user_id"

	// EventGetTokenFromStore is emitted before fetching the token from the cache.
	EventGetTokenFromStore = "get_token_from_store"

	// EventCompareToken is emitted before comparing the provided token with the stored one.
	EventCompareToken = "compare_token"

	// EventTokenValidated is emitted after successful validation.
	EventTokenValidated = "token_validated"
)
