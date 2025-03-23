package http

import (
	"github.com/lechitz/AionApi/core/domain"
)

type IUserService interface {
	CreateUser(ctx domain.ContextControl, userDomain domain.UserDomain, password string) (domain.UserDomain, error)
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(ctx domain.ContextControl, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error)
	UpdateUser(ctx domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx domain.ContextControl, userDomain domain.UserDomain, password, newPassword string) (domain.UserDomain, string, error)
	SoftDeleteUser(ctx domain.ContextControl, ID uint64) error
	HashPassword(password string) (string, error)
}
