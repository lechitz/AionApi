package msg

// Error Messages
const (
	ErrorToStoreTokenInRedis      = "failed to store token in Redis"
	ErrorToDeleteTokenFromRedis   = "error to delete token from Redis"
	ErrorToRetrieveTokenFromRedis = "error retrieving token from Redis"
	ErrorToUpdateToken            = "error updating token"
)

// Error Messages to Validate Token
const (
	ErrorInvalidToken       = "invalid token"
	ErrorInvalidTokenClaims = "invalid token claims"
	ErrorInvalidUserIDClaim = "invalid userID claim in token"
)

// Success Messages
const (
	SuccessTokenCreated   = "token created successfully"
	SuccessTokenValidated = "token validated successfully"
	InfoTokenRetrieved    = "token retrieved successfully"
	SuccessTokenDeleted   = "token deleted successfully"
	SuccessTokenUpdated   = "token updated successfully"
)
