// Package domain contains core business entities used throughout the application.
package domain

// Auth represents a token associated with a user in the system.
// It includes the token string and the corresponding user ID.
//
// The Auth struct is used for both access and refresh tokens. The 'Type'
// field (recommended values are TokenTypeAccess and TokenTypeRefresh) helps
// distinguish token semantics at runtime when needed.
type Auth struct {
	Token string // JWT or session token string
	Key   uint64 // ID of the user to whom the token belongs
	Type  string // "access" or "refresh"
}

// RefreshAuth represents a refresh token associated with a user.
// Kept for backward compatibility; prefer using RefreshToken type and
// the constructors below when creating tokens in new code.
//
// NOTE: RefreshAuth is a separate struct historically used by older code.
// It has the same fields as Auth except for the optional Type field. New
// code should prefer AccessToken / RefreshToken so the Type field is
// consistently available.
type RefreshAuth struct {
	Token string
	Key   uint64
}

// Token type constants. Use these values when setting or checking the
// Auth.Type field to avoid string literals scattered through the codebase.
const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

// AccessToken is a typed representation of an access token based on Auth.
// It is a distinct type at compile time to prevent accidental interchange
// with RefreshToken without explicit conversion.
type AccessToken Auth

// RefreshToken is a typed representation of a refresh token based on Auth.
// It is a distinct type at compile time to prevent accidental interchange
// with AccessToken without explicit conversion.
type RefreshToken Auth

// NewAccessToken creates an AccessToken with the canonical Type value set.
// Use this helper to ensure the Type field is consistently populated.
func NewAccessToken(token string, key uint64) AccessToken {
	return AccessToken{Token: token, Key: key, Type: TokenTypeAccess}
}

// NewRefreshToken creates a RefreshToken with the canonical Type value set.
func NewRefreshToken(token string, key uint64) RefreshToken {
	return RefreshToken{Token: token, Key: key, Type: TokenTypeRefresh}
}

// ToAuth converts an AccessToken to the underlying Auth type.
func (a AccessToken) ToAuth() Auth { return Auth(a) }

// ToAuth converts a RefreshToken to the underlying Auth type.
func (r RefreshToken) ToAuth() Auth { return Auth(r) }

// AccessTokenFromAuth converts an Auth to an AccessToken.
func AccessTokenFromAuth(a Auth) AccessToken { return AccessToken(a) }

// RefreshTokenFromAuth converts an Auth to a RefreshToken.
func RefreshTokenFromAuth(a Auth) RefreshToken { return RefreshToken(a) }
