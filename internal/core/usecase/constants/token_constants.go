package constants

import "time"

const (
	ErrorToSaveToken   = "error to save token"
	ErrorToUpdateToken = "error to update token"
	ErrorToDeleteToken = "error to delete token"
	ErrorToAssignToken = "error to assign token"

	ErrorInvalidToken       = "invalid token"
	ErrorInvalidTokenClaims = "invalid token claims"
	ErrorInvalidUserIDClaim = "invalid userID in claim"

	ErrorToRetrieveTokenFromCache = "error to retrieve token from cache"
	ErrorTokenMismatch            = "provided token does not match stored token"

	SuccessTokenCreated   = "token created successfully"
	SuccessTokenValidated = "token validated successfully"
	SuccessTokenUpdated   = "token updated successfully"
	SuccessTokenDeleted   = "token deleted successfully"

	Token           = "token"
	Exp             = "exp"
	ExpTimeToken    = 1 * time.Hour
	TokenFromCookie = "TokenFromCookie"
	TokenFromCache  = "TokenFromCache"

	SecretKey = "secret"
)
