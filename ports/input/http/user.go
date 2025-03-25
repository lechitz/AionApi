package http

import (
	"github.com/lechitz/AionApi/core/domain"
)

type IUserService interface {
	CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error)
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error)
	UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
	SoftDeleteUser(ctx domain.ContextControl, id uint64) error
}
