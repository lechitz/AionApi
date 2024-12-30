package input

import "github.com/lechitz/AionApi/internal/core/domain"

type IUserService interface {
	CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
}
