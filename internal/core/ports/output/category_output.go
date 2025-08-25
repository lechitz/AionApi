// Package output db defines interfaces for managing categories in the system.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

// CategoryCreator defines an interface for creating a new category with specified attributes in a given context.
type CategoryCreator interface {
	Create(ctx context.Context, category domain.Category) (domain.Category, error)
}

// CategoryRetriever defines methods for retrieving category details or multiple categories for a specific user.
type CategoryRetriever interface {
	GetByID(ctx context.Context, categoryID, userID uint64) (domain.Category, error)
	GetByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error)
	ListAll(ctx context.Context, userID uint64) ([]domain.Category, error)
}

// CategoryUpdater defines an interface for updating a category's attributes in a given context based on the provided category ID and user ID.
type CategoryUpdater interface {
	UpdateCategory(ctx context.Context, categoryID uint64, userID uint64, fields map[string]interface{}) (domain.Category, error)
}

// CategoryDeleter defines an interface for handling the soft deletion of categories within a contextual operation.
type CategoryDeleter interface {
	SoftDelete(ctx context.Context, categoryID, userID uint64) error
}

// CategoryRepository represents a composite interface for managing categories, combining creation, retrieval, updating, and soft-deletion functionalities.
type CategoryRepository interface {
	CategoryCreator
	CategoryRetriever
	CategoryUpdater
	CategoryDeleter
}
