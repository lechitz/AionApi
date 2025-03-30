package http

import (
	"github.com/lechitz/AionApi/internal/core/domain"
)

type IAuthService interface {
	Login(ctx domain.ContextControl, userDomain domain.UserDomain, passwordReq string) (domain.UserDomain, string, error)
	Logout(ctx domain.ContextControl, token string) error
}
