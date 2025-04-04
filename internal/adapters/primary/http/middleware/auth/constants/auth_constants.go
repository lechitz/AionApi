package constants

// Error messages for the auth middleware
const (
	ErrorUnauthorizedAccessMissingToken = "Unauthorized access: missing token"
	ErrorUnauthorizedAccessInvalidToken = "Unauthorized access: invalid token"
)

// Success messages for the auth middleware

const (
	SuccessTokenValidated = "Token validated successfully, adding userID to context"
)

const (
	Error     = "error"
	UserID    = "id"
	Token     = "token"
	AuthToken = "auth_token"
)
