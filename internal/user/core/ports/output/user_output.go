// Package output defines interfaces for user-related output ports.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/user/core/domain"
)

// UserCreator defines an interface for creating a new user in the system.
type UserCreator interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

// UserUniqueness holds information about whether a username or email is already taken
// and, if so, who owns them.
type UserUniqueness struct {
	UsernameOwnerID *uint64
	EmailOwnerID    *uint64
	UsernameTaken   bool
	EmailTaken      bool
}

// UserFinder defines method for retrieving a user by ID.
type UserFinder interface {
	CheckUniqueness(ctx context.Context, username, email string) (UserUniqueness, error)
	GetByID(ctx context.Context, userID uint64) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	ListAll(ctx context.Context) ([]domain.User, error)
}

// UserUpdater defines an interface for updating user information in the system.
type UserUpdater interface {
	Update(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.User, error)
}

// UserDeleter defines an interface for handling user deletion within the system.
type UserDeleter interface {
	SoftDelete(ctx context.Context, userID uint64) error
}

// UserRepository aggregates interfaces for managing user creation, retrieval, updating, and deletion in the system.
type UserRepository interface {
	UserCreator
	UserFinder
	UserUpdater
	UserDeleter
}
