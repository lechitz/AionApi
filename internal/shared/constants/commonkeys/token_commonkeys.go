// Package commonkeys defines shared keys for token fields used in config, logs, and context.
package commonkeys

const (
	// Token is the key for identifying a generic token value in configs or logs.
	Token = "token"

	// AuthTokenCookieName is the key for an authentication token value.
	AuthTokenCookieName = "auth_token"

	// TokenFromCookie is the key for a token value retrieved from cookies.
	TokenFromCookie = "token_from_cookie"

	// TokenFromCache is the key for a token value retrieved from cache.
	TokenFromCache = "token_from_cache"
)
