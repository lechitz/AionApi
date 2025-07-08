package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// UserCreator defines an interface for creating a new user in the system.
type UserCreator interface {
	CreateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
}

// UserRetriever defines methods for retrieving user data by various identifiers.
type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
	GetUserByID(ctx context.Context, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error)
}

// UserUpdater defines an interface for updating user information in the system.
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.UserDomain, error)
}

// UserDeleter defines an interface for handling user deletion within the system.
type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, userID uint64) error
}

// UserStore aggregates interfaces for managing user creation, retrieval, updating, and deletion in the system.
type UserStore interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
