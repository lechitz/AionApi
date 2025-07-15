package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// UserCreator defines an interface for creating a new user in the system.
type UserCreator interface {
	CreateUser(ctx context.Context, userDomain domain.User) (domain.User, error)
}

//TODO:"deixar somente um método para recuperar o usuário indformando o parametro?! Pensar..

// UserByIDFinder defines method for retrieving a user by ID.
type UserByIDFinder interface {
	GetUserByID(ctx context.Context, userID uint64) (domain.User, error)
}

// UserByUsernameFinder defines method for retrieving a user by username.
type UserByUsernameFinder interface {
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

// UserByEmailFinder defines method for retrieving a user by email.
type UserByEmailFinder interface {
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

// UserLister defines method for listing all users.
type UserLister interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
}

// TODO: deletar a interface UserRetriever depois que terminar o uso dela !!!

// UserRetriever defines methods for retrieving user data by various identifiers.
type UserRetriever interface {
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, ID uint64) (domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

// TODO: *

// UserUpdater defines an interface for updating user information in the system.
type UserUpdater interface {
	UpdateUser(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.User, error)
}

// UserDeleter defines an interface for handling user deletion within the system.
type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, userID uint64) error
}

// UserRepository aggregates interfaces for managing user creation, retrieval, updating, and deletion in the system.
type UserRepository interface {
	UserCreator
	UserRetriever
	UserUpdater
	UserDeleter
}
