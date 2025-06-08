package constants

// Error messages related to token operations.
const (
	ErrorToSaveToken              = "error saving access reference"
	ErrorToUpdateToken            = "error updating access reference"
	ErrorToDeleteToken            = "error deleting access reference"
	ErrorToAssignToken            = "error assigning access reference"
	ErrorInvalidToken             = "invalid access reference"
	ErrorInvalidTokenClaims       = "invalid claims in reference"
	ErrorInvalidUserIDClaim       = "invalid userID in claims"
	ErrorToRetrieveTokenFromCache = "error retrieving access reference from cache"
	ErrorTokenMismatch            = "provided reference does not match stored one"
)

// Success messages related to token lifecycle operations.
const (
	SuccessTokenCreated   = "access reference created successfully"
	SuccessTokenValidated = "access reference validated successfully"
	SuccessTokenUpdated   = "access reference updated successfully"
	SuccessTokenDeleted   = "access reference deleted successfully"
)

// Token-related constant keys used across application contexts.
const (
	Token           = "token"
	TokenFromCookie = "TokenFromCookie"
	TokenFromCache  = "TokenFromCache"
	SecretKey       = "secret"
	Error           = "error"
	UserID          = "user_id"
)
