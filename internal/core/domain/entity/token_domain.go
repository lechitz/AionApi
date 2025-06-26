// Package entity contains core business entities used throughout the application.
package entity

// TokenConfig holds configuration data related to token handling,
// such as the secret key used for signing JWTs.
type TokenConfig struct {
	SecretKey string // Secret key for signing/verifying tokens
}

// TokenDomain represents a token associated with a user in the system.
// It includes the token string and the corresponding user ID.
type TokenDomain struct {
	Token  string // JWT or session token string
	UserID uint64 // ID of the user to whom the token belongs
}
