package http

import "github.com/lechitz/AionApi/internal/core/domain"

type UserCreator interface {
	CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error)
}

type UserRetriever interface {
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
	GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error)
}

type UserUpdater interface {
	UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
}

type UserDeleter interface {
	SoftDeleteUser(ctx domain.ContextControl, id uint64) error
}

type UserService interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
