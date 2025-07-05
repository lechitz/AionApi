// Package input http defines interfaces for user authentication and session management.
package input

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

// Authenticator defines methods for authenticating users and generating access tokens.
type Authenticator interface {
	// Login validates user credentials and returns the user, token, and any error encountered.
	Login(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, string, error)
}

// SessionRevoker provides a method to invalidate user sessions by revoking tokens.
type SessionRevoker interface {
	// Logout revokes a session based on the supplied token and context, returning an error if unsuccessful.
	Logout(ctx context.Context, token string) error
}

// AuthService combines the functionalities of authentication and session management.
// It includes methods for user login, token generation, and session revocation.
type AuthService interface {
	Authenticator
	SessionRevoker
}
