// Package commonkeys defines shared keys for token fields used in config, logs, and context.
package commonkeys

const (
	// Token is the key for identifying a generic token value in configs or logs.
	Token = "token"

	// TokenKey is the key for identifying a token key.
	TokenKey = "token_key"

	// AuthTokenCookieName is the key for an authentication token value.
	AuthTokenCookieName = "auth_token"
)
