package input

import "github.com/lechitz/AionApi/internal/core/domain"

type IAuthService interface {
	Login(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, string, error)
	Logout(contextControl domain.ContextControl, userDomain domain.UserDomain, token string) error
}
