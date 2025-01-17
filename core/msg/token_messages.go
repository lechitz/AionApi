package msg

// Error Messages

const (
	ErrorToCreateToken            = "error creating authentication token"
	ErrorTokenMismatch            = "provided token does not match stored token"
	ErrorRevokeToken              = "failed to revoke token during logout"
	ErrorToUpdatePassword         = "error updating password"
	ErrorToDeleteToken            = "error deleting token"
	ErrorToAssignToken            = "error assigning token"
	ErrorToGenerateToken          = "error generating token"
	ErrorToSaveToken              = "error saving token"
	ErrorToCheckToken             = "error to check token"
	ErrorToRetrieveTokenFromCache = "error retrieving token from cache"
	ErrorToUpdateToken            = "error updating token"
	ErrorInvalidToken             = "invalid token"
	ErrorInvalidTokenClaims       = "invalid token claims"
	ErrorInvalidUserIDClaim       = "invalid userID claim in token"
)

// Success Messages

const (
	SuccessTokenCreated   = "token created successfully"
	SuccessTokenValidated = "token validated successfully"
	SuccessTokenDeleted   = "token deleted successfully"
	SuccessTokenUpdated   = "token updated successfully"
)

// Messages

const (
	TokenFromCookie = "tokenFromCookie"
	TokenFromCache  = "tokenFromCache"
)
