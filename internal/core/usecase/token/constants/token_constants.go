package constants

// Error messages related to token operations.
const (
	// ErrorToSaveToken indicates failure while saving a token.
	ErrorToSaveToken = "error to save token"

	// ErrorToUpdateToken indicates failure while updating a token.
	ErrorToUpdateToken = "error to update token"

	// ErrorToDeleteToken indicates failure while deleting a token.
	ErrorToDeleteToken = "error to delete token"

	// ErrorToAssignToken indicates failure while assigning a token to a user/session.
	ErrorToAssignToken = "error to assign token"

	// ErrorInvalidToken is used when a token is malformed or expired.
	ErrorInvalidToken = "invalid token"

	// ErrorInvalidTokenClaims is returned when token claims are missing or incorrect.
	ErrorInvalidTokenClaims = "invalid token claims"

	// ErrorInvalidUserIDClaim is returned when the 'user_id' claim is invalid or missing.
	ErrorInvalidUserIDClaim = "invalid userID in claim"

	// ErrorToRetrieveTokenFromCache indicates a failure while retrieving the token from cache.
	ErrorToRetrieveTokenFromCache = "error to retrieve token from cache"

	// ErrorTokenMismatch is returned when the provided token doesn't match the cached token.
	ErrorTokenMismatch = "provided token does not match stored token"
)

// Success messages related to token lifecycle operations.
const (
	// SuccessTokenCreated indicates that a token was successfully created.
	SuccessTokenCreated = "token created successfully"

	// SuccessTokenValidated indicates that a token was successfully validated.
	SuccessTokenValidated = "token validated successfully"

	// SuccessTokenUpdated indicates that a token was successfully updated.
	SuccessTokenUpdated = "token updated successfully"

	// SuccessTokenDeleted indicates that a token was successfully deleted.
	SuccessTokenDeleted = "token deleted successfully"
)

// Token-related constant keys used across application contexts.
const (
	// Token is the generic key for accessing token values in contexts.
	Token = "token"

	// TokenFromCookie is the context key for the token extracted from cookies.
	TokenFromCookie = "TokenFromCookie"

	// TokenFromCache is the context key for the token extracted from Redis or another cache.
	TokenFromCache = "TokenFromCache"

	// SecretKey represents the environment variable or config field name for the JWT secret.
	SecretKey = "secret"

	// Error is a generic error key used in response maps or logging.
	Error = "error"

	// UserID is the key used to reference a user ID inside tokens or request contexts.
	UserID = "user_id"
)
