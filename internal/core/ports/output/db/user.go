package db

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
}

type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
	GetUserByID(ctx context.Context, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error)
}

type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.UserDomain, error)
}

type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, userID uint64) error
}

type UserStore interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
