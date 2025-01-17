package msg

// Error Messages
const (
	ErrorToCompareHashAndPassword = "invalid password provided"
	ErrorToCreateToken            = "error creating authentication token"
	ErrorInvalidUserID            = "invalid userID provided for logout"
	ErrorTokenNotFound            = "token not found for user in Redis"
	ErrorTokenMismatch            = "provided token does not match stored token"
	ErrorRetrieveTokenFromRedis   = "error retrieving token from Redis"
	ErrorRevokeToken              = "failed to revoke token during logout"
	ErrorToHashPassword           = "error hashing password"
	ErrorToUpdatePassword         = "error updating password"
	ErrorToDeleteToken            = "error deleting token"
	ErrorToRevokeToken            = "error revoking token"
	ErrorToAssignToken            = "error assigning token"
	ErrorToGenerateToken          = "error generating token"
	ErrorToSaveToken              = "error saving token"
	ErrorToCheckToken             = "error to check token"
)

// Success Messages
const (
	SuccessUserLoggedOut   = "user logged out successfully"
	SuccessPasswordUpdated = "password updated successfully"
	SuccessTokenRetrieved  = "token retrieved successfully"

)

const (
	StringTokenKey = "user:%d:token"
)
