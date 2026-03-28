// Package input user defines interfaces for user management and authentication.
package input

import (
	"context"

	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// UserCreator defines a method to create a new user with the given details and return the created user or an error.
type UserCreator interface {
	Create(ctx context.Context, cmd CreateUserCommand) (domain.User, error)
}

// UserRegistrationFlow defines the multi-step public registration process.
type UserRegistrationFlow interface {
	StartRegistration(ctx context.Context, cmd StartRegistrationCommand) (domain.RegistrationSession, error)
	UpdateRegistrationProfile(ctx context.Context, registrationID string, cmd UpdateRegistrationProfileCommand) (domain.RegistrationSession, error)
	UpdateRegistrationAvatar(ctx context.Context, registrationID string, cmd UpdateRegistrationAvatarCommand) (domain.RegistrationSession, error)
	CompleteRegistration(ctx context.Context, registrationID string) (domain.User, error)
}

// UserAvatarUploader defines a method to validate/process uploaded avatar and return a usable avatar URL.
type UserAvatarUploader interface {
	UploadAvatar(ctx context.Context, cmd UploadAvatarCommand) (string, string, int64, error)
}

// UserReader defines methods for retrieving user information by unique identifier.
type UserReader interface {
	GetByID(ctx context.Context, userID uint64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
	GetUserStats(ctx context.Context, userID uint64) (domain.UserStats, error)
}

// UserUpdater defines methods for updating user information and user passwords in the system.
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, cmd UpdateUserCommand) (domain.User, error)
	RemoveAvatar(ctx context.Context, userID uint64) (domain.User, error)
	UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error)
}

// UserRemover provides a method to perform a soft delete operation on a user by their unique identifier.
type UserRemover interface {
	SoftDeleteUser(ctx context.Context, userID uint64) error
}

// UserService defines methods for managing user creation, retrieval, updates, and soft deletion by extending related interfaces.
type UserService interface {
	UserCreator
	UserRegistrationFlow
	UserAvatarUploader
	UserReader
	UserUpdater
	UserRemover
}
