// Package db defines interfaces for managing categories in the system.
package db

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
)

// CategoryCreator defines an interface for creating a new category with specified attributes in a given context.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error)
}

// CategoryRetriver defines methods for retrieving category details or multiple categories for a specific user.
// GetCategoryByID fetches a single category by its unique identifier asynchronously.
// GetCategoryByName retrieves a category based on its name from the system.
// GetAllCategories returns a list of all categories associated with a specific user ID.
type CategoryRetriver interface {
	GetCategoryByID(ctx context.Context, category entity.Category) (entity.Category, error)
	GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]entity.Category, error)
}

// CategoryUpdater defines an interface for updating a category's attributes in a given context based on the provided category ID and user ID.
type CategoryUpdater interface {
	UpdateCategory(
		ctx context.Context,
		categoryID uint64,
		userID uint64,
		fields map[string]interface{},
	) (entity.Category, error)
}

// CategoryDeleter defines an interface for handling the soft deletion of categories within a contextual operation.
type CategoryDeleter interface {
	SoftDeleteCategory(ctx context.Context, category entity.Category) error
}

// CategoryStore represents a composite interface for managing categories, combining creation, retrieval, updating, and soft-deletion functionalities.
type CategoryStore interface {
	CategoryCreator
	CategoryRetriver
	CategoryUpdater
	CategoryDeleter
}
