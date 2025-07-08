// Package input graphql defines a contract for interacting with category-related operations.
package input

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// CategoryCreator defines a contract for creating a new category within the system.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

// CategoryRetriver defines methods for retrieving category data based on criteria or user associations.
type CategoryRetriver interface {
	GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error)
}

// CategoryUpdater defines a contract for updating an existing category within the system.
type CategoryUpdater interface {
	UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

// CategoryDeleter defines a contract for soft-deleting a category without permanently removing it from the system.
type CategoryDeleter interface {
	SoftDeleteCategory(ctx context.Context, category domain.Category) error
}

// CategoryService defines a contract that combines creating, retrieving, updating, and soft-deleting categories in the system.
type CategoryService interface {
	CategoryCreator
	CategoryRetriver
	CategoryUpdater
	CategoryDeleter
}
