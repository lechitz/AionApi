package constants

// Error Messages
const (
	ErrorToStoreTokenInRedis      = "Failed to store token in Redis"
	ErrorToDeleteTokenFromRedis   = "Error to delete token from Redis"
	ErrorToRetrieveTokenFromRedis = "Error retrieving token from Redis"
	ErrorToUpdateToken            = "Error updating token"
)

// Error Messages to Validate Token
const (
	ErrorInvalidToken       = "Invalid token"
	ErrorInvalidTokenClaims = "Invalid token claims"
	ErrorInvalidUserIDClaim = "Invalid userID claim in token"
)

// Success Messages
const (
	SuccessTokenCreated   = "Token created successfully"
	SuccessTokenValidated = "Token validated successfully"
	InfoTokenRetrieved    = "Token retrieved successfully"
	SuccessTokenDeleted   = "Token deleted successfully"
	SuccessTokenUpdated   = "Token updated successfully"
)
