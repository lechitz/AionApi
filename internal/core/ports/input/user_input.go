// Package input user defines interfaces for user management and authentication.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// UserCreator defines a method to create a new user with the given details and return the created user or an error.
type UserCreator interface {
	CreateUser(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, error)
}

// UserRetriever defines methods to retrieve user information from a data source.
type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
	GetUserByID(ctx context.Context, id uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
}

// UserUpdater defines methods for updating user information and user passwords in the system.
type UserUpdater interface {
	UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
}

// UserDeleter provides a method to perform a soft delete operation on a user by their unique identifier.
type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, id uint64) error
}

// UserService defines methods for managing user creation, retrieval, updates, and soft deletion by extending related interfaces.
type UserService interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
