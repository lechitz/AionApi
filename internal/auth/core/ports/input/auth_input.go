// Package input http defines interfaces for user authentication and session management.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// Validator validates a raw token and returns the resolved userID and claims.
type Validator interface {
	Validate(ctx context.Context, token string) (userID uint64, claims map[string]any, err error)
}

// Authenticator defines methods for authentication users and generating access tokens.
type Authenticator interface {
	Login(ctx context.Context, username, password string) (domain.User, string, string, error)
}

// SessionRevoker provides a method to invalidate user sessions by revoking tokens.
type SessionRevoker interface {
	Logout(ctx context.Context, userID uint64) error
}

// Refresher provides a method to renew access tokens using a refresh token.
type Refresher interface {
	RefreshTokenRenewal(ctx context.Context, refreshToken string) (newAccessToken string, newRefreshToken string, err error)
}

// AuthService combines the functionalities of authentication and session management.
type AuthService interface {
	Authenticator
	Validator
	SessionRevoker
	Refresher
}
