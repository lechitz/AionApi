package http

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, error)
}

type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
	GetUserByID(ctx context.Context, id uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
}

type UserUpdater interface {
	UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
}

type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, id uint64) error
}

type UserService interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
