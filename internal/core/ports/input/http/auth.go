package http

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

type Authenticator interface {
	Login(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, string, error)
}

type SessionRevoker interface {
	Logout(ctx domain.ContextControl, token string) error
}

type AuthService interface {
	Authenticator
	SessionRevoker
}
