package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
)

// UserCreator defines a method to create a new user with the given details and return the created user or an error.
type UserCreator interface {
	CreateUser(ctx context.Context, user entity.UserDomain, password string) (entity.UserDomain, error)
}

// UserRetriever defines methods to retrieve user information from a data source.
// GetAllUsers fetches all users.
// GetUserByID retrieves a user by their unique identifier.
// GetUserByUsername retrieves a user by their unique username.
type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]entity.UserDomain, error)
	GetUserByID(ctx context.Context, id uint64) (entity.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (entity.UserDomain, error)
}

// UserUpdater defines methods for updating user information and user passwords in the system.
// UpdateUser updates user details in the data source and returns the updated user or an error.
// UpdateUserPassword updates the user's password and returns the updated user, a result message, or an error.
type UserUpdater interface {
	UpdateUser(ctx context.Context, user entity.UserDomain) (entity.UserDomain, error)
	UpdateUserPassword(ctx context.Context, user entity.UserDomain, oldPassword, newPassword string) (entity.UserDomain, string, error)
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
