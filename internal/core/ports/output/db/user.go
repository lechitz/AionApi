package db

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// UserCreator defines an interface for creating a new user in the system.
// It requires a context and a user domain object as input and returns the created user or an error.
type UserCreator interface {
	CreateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
}

// UserRetriever defines methods for retrieving user data by various identifiers.
// It supports operations for fetching all users, or by ID, username, or email.
type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
	GetUserByID(ctx context.Context, ID uint64) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error)
}

// UserUpdater defines an interface for updating user information in the system.
// It accepts a context, user ID, and a map of fields to update, returning the updated user or an error.
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.UserDomain, error)
}

// UserDeleter defines an interface for handling user deletion within the system.
// It enables soft deletion of users by their unique ID in a given context.
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
