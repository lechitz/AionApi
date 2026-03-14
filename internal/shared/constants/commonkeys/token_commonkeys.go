// Package commonkeys defines shared keys for token fields used in config, logs, and context.
package commonkeys

const (
	// Token is the key for identifying a generic token value in configs or logs.
	Token = "token"

	// TokenKey is the key for identifying a token key.
	TokenKey = "token_key"

	// AuthTokenCookieName is the key for an authentication token value.
	AuthTokenCookieName = "auth_token"

	// GraceKey is the log/context key used for grace period redis keys.
	GraceKey = "grace_key"

	// GraceTTL is the log/context key used for grace period TTL values.
	GraceTTL = "grace_ttl"

	// OldTokenPrefix is the log/context key for old token prefix.
	OldTokenPrefix = "old_token_prefix"

	// NewTokenPrefix is the log/context key for new token prefix.
	NewTokenPrefix = "new_token_prefix"

	// ProvidedTokenPrefix is the log/context key for provided token prefix (validation).
	ProvidedTokenPrefix = "provided_prefix"

	// CachedTokenPrefix is the log/context key for cached token prefix (validation).
	CachedTokenPrefix = "cached_prefix" // #nosec G101: log key, not a credential

	// GraceLookupError is the log/context key for grace period lookup errors.
	GraceLookupError = "grace_lookup_error"
)

// Token type constants.
const (
	// TokenTypeAccess identifies an access token stored in the cache.
	TokenTypeAccess = "access"

	// TokenTypeRefresh identifies a refresh token stored in the cache.
	TokenTypeRefresh = "refresh"
)
