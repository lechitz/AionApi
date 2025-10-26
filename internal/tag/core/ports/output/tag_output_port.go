// Package output db defines interfaces for managing tags in the system.
package output

import (
	"context"

	"github.com/lechitz/AionApi/internal/tag/core/domain"
)

// TagCreator defines an interface for creating a new tag with specified attributes in a given context.
type TagCreator interface {
	Create(ctx context.Context, tag domain.Tag) (domain.Tag, error)
}

// TagRetriever defines methods for retrieving tag details or multiple tags for a specific user.
type TagRetriever interface {
	GetByID(ctx context.Context, tagID, userID uint64) (domain.Tag, error)
	GetByName(ctx context.Context, TagName string, userID uint64) (domain.Tag, error)
	GetByCategoryID(ctx context.Context, categoryID uint64, userID uint64) ([]domain.Tag, error)
	GetAll(ctx context.Context, userID uint64) ([]domain.Tag, error)
}

// TagRepository represents a composite interface for managing tags, combining creation, retrieval, updating, and soft-deletion functionalities.
type TagRepository interface {
	TagCreator
	TagRetriever
}
