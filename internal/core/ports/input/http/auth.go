package http

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

type Authenticator interface {
	Login(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, string, error)
}

type SessionRevoker interface {
	Logout(ctx context.Context, token string) error
}

type AuthService interface {
	Authenticator
	SessionRevoker
}
