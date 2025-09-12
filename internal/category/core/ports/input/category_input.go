// Package input graphql defines a contract for interacting with handler-related operations.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/category/core/domain"
)

// CategoryCreator defines a contract for creating a new handler within the system.
type CategoryCreator interface {
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
}

// CategoryRetriever defines methods for retrieving handler data based on criteria or user associations.
type CategoryRetriever interface {
	GetByID(ctx context.Context, categoryID, userID uint64) (domain.Category, error)
	GetByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error)
	ListAll(ctx context.Context, userID uint64) ([]domain.Category, error)
}

// CategoryUpdater defines a contract for updating an existing handler within the system.
type CategoryUpdater interface {
	Update(ctx context.Context, category domain.Category) (domain.Category, error)
}

// CategoryDeleter defines a contract for soft-deleting a handler without permanently removing it from the system.
type CategoryDeleter interface {
	SoftDelete(ctx context.Context, categoryID, userID uint64) error
}

// CategoryService defines a contract that combines creating, retrieving, updating, and soft-deleting categories in the system.
type CategoryService interface {
	CategoryCreator
	CategoryRetriever
	CategoryUpdater
	CategoryDeleter
}
