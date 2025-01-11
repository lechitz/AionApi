package constants

// Error Messages
const (
	ErrorToVerifyPassword       = "invalid password provided"
	ErrorToCreateToken          = "error creating authentication token"
	ErrorInvalidUserID          = "invalid userID provided for logout"
	ErrorTokenNotFound          = "token not found for user in Redis"
	ErrorTokenMismatch          = "provided token does not match stored token"
	ErrorRetrieveTokenFromRedis = "error retrieving token from Redis"
	ErrorRevokeToken            = "failed to revoke token during logout"
	ErrorToHashPassword         = "error hashing password"
	ErrorToUpdatePassword       = "error updating password"
	ErrorToDeleteToken          = "error deleting token"
)

// Success Messages
const (
	SuccessUserLoggedOut   = "user logged out successfully"
	SuccessPasswordUpdated = "password updated successfully"
)
