// Package input user defines interfaces for user management and authentication.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// CreateUserCommand defines a struct for creating a new user.
type CreateUserCommand struct {
	Name     string
	Username string
	Email    string
	Password string
}

// UserCreator defines a method to create a new user with the given details and return the created user or an error.
type UserCreator interface {
	Create(ctx context.Context, cmd CreateUserCommand) (domain.User, error)
}

type UserReader interface {
	GetByID(ctx context.Context, userID uint64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
}

// UpdateUserCommand defines a struct for updating user information.
type UpdateUserCommand struct {
	Name     *string
	Username *string
	Email    *string
}

func (u UpdateUserCommand) HasUpdates() bool {
	return u.Name != nil || u.Username != nil || u.Email != nil
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
