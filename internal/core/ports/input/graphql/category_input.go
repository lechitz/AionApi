// Package graphql defines a contract for interacting with category-related operations.
package graphql

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
)

// CategoryCreator defines a contract for creating a new category within the system.
// It encapsulates logic to persist category data and ensure consistency.
// CreateCategory creates a category with given metadata and returns it or an error.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
}

// CategoryRetriver defines methods for retrieving category data based on criteria or user associations.
// GetCategoryByID retrieves a specific category by its ID.
// GetCategoryByName retrieves a specific category by its name.
// GetAllCategories retrieves all categories associated with a specific user ID.
type CategoryRetriver interface {
	GetCategoryByID(ctx context.Context, category entity.Category) (entity.Category, error)
	GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]entity.Category, error)
}

// CategoryUpdater defines a contract for updating an existing category within the system.
// UpdateCategory updates category metadata and returns the updated entity or an error.
type CategoryUpdater interface {
	UpdateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
}

// CategoryDeleter defines a contract for soft-deleting a category without permanently removing it from the system.
type CategoryDeleter interface {
	SoftDeleteCategory(ctx context.Context, category entity.Category) error
}

// CategoryService defines a contract that combines creating, retrieving, updating, and soft-deleting categories in the system.
type CategoryService interface {
	CategoryCreator
	CategoryRetriver
	CategoryUpdater
	CategoryDeleter
}
