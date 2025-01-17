package http

import (
	"github.com/lechitz/AionApi/core/domain/entities"
)

type IUserService interface {
	CreateUser(ctx entities.ContextControl, userDomain entities.UserDomain, password string) (entities.UserDomain, error)
	GetAllUsers(ctx entities.ContextControl) ([]entities.UserDomain, error)
	GetUserByID(ctx entities.ContextControl, ID uint64) (entities.UserDomain, error)
	GetUserByUsername(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error)
	UpdateUser(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error)
	UpdateUserPassword(ctx entities.ContextControl, userDomain entities.UserDomain, password, newPassword string) (entities.UserDomain, string, error)
	SoftDeleteUser(ctx entities.ContextControl, ID uint64) error
	HashPassword(password string) (string, error)
}
