package db

import (
	"github.com/lechitz/AionApi/core/domain/entities"
)

type IUserRepository interface {
	CreateUser(ctx entities.ContextControl, user entities.UserDomain) (entities.UserDomain, error)
	GetAllUsers(ctx entities.ContextControl) ([]entities.UserDomain, error)
	GetUserByID(ctx entities.ContextControl, ID uint64) (entities.UserDomain, error)
	GetUserByUsername(ctx entities.ContextControl, username string) (entities.UserDomain, error)
	UpdateUser(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error)
	UpdatePassword(ctx entities.ContextControl, userDomain entities.UserDomain) (entities.UserDomain, error)
	SoftDeleteUser(ctx entities.ContextControl, ID uint64) error
}
