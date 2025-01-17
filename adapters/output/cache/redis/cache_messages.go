package redis

// Errors messages related to token
const (
	ErrorToCreateToken            = "error creating authentication token"
	ErrorToStoreTokenInRedis      = "error storing token in Redis"
	ErrorInvalidToken             = "invalid token"
	ErrorInvalidTokenClaims       = "invalid token claims"
	ErrorInvalidUserIDClaim       = "invalid userID claim in token"
	ErrorToRetrieveTokenFromRedis = "error retrieving token from Redis"
	ErrorTokenMismatch            = "provided token does not match stored token"
	ErrorToUpdateToken            = "error updating token"
	ErrorToDeleteTokenFromRedis   = "error deleting token from Redis"
)

// Success messages related to token
const (
	SuccessTokenCreated   = "token created successfully"
	SuccessTokenValidated = "token validated successfully"
	SuccessTokenRetrieved = "token retrieved successfully"
	SuccessTokenUpdated   = "token updated successfully"
	SuccessTokenDeleted   = "token deleted successfully"
)
