package input

import "github.com/lechitz/AionApi/internal/core/domain"

type IUserService interface {
	CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(contextControl domain.ContextControl, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(contextControl domain.ContextControl, username domain.UserDomain) (domain.UserDomain, error)
	UpdateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	SoftDeleteUser(contextControl domain.ContextControl, ID uint64) error
}
