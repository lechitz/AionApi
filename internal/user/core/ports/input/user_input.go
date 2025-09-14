// Package input user defines interfaces for user management and authentication.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// UserCreator defines a method to create a new user with the given details and return the created user or an error.
type UserCreator interface {
	Create(ctx context.Context, cmd CreateUserCommand) (domain.User, error)
}

// UserReader defines methods for retrieving user information by unique identifier.
type UserReader interface {
	GetByID(ctx context.Context, userID uint64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
}

// UserUpdater defines methods for updating user information and user passwords in the system.
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, cmd UpdateUserCommand) (domain.User, error)
	UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error)
}

// UserRemover provides a method to perform a soft delete operation on a user by their unique identifier.
type UserRemover interface {
	SoftDeleteUser(ctx context.Context, userID uint64) error
}

// UserService defines methods for managing user creation, retrieval, updates, and soft deletion by extending related interfaces.
type UserService interface {
	UserCreator
	UserReader
	UserUpdater
	UserRemover
}
