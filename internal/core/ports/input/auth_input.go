// Package input http defines interfaces for user authentication and session management.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// Authenticator defines methods for authenticating users and generating access tokens.
type Authenticator interface {
	Login(ctx context.Context, username, password string) (domain.User, string, error)
}

// SessionRevoker provides a method to invalidate user sessions by revoking tokens.
type SessionRevoker interface {
	Logout(ctx context.Context, userID uint64) error
}

// AuthService combines the functionalities of authentication and session management.
type AuthService interface {
	Authenticator
	SessionRevoker
}
