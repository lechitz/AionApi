// Package output defines interfaces for category-related cache operations.
package output

import (
	"context"
	"time"

	"github.com/lechitz/AionApi/internal/category/core/domain"
)

// CategoryCache defines cache operations for categories.
type CategoryCache interface {
	SaveCategory(ctx context.Context, category domain.Category, expiration time.Duration) error
	SaveCategoryByName(ctx context.Context, category domain.Category, expiration time.Duration) error
	SaveCategoryList(ctx context.Context, userID uint64, categories []domain.Category, expiration time.Duration) error
	GetCategory(ctx context.Context, categoryID, userID uint64) (domain.Category, error)
	GetCategoryByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error)
	GetCategoryList(ctx context.Context, userID uint64) ([]domain.Category, error)
	DeleteCategory(ctx context.Context, categoryID, userID uint64) error
	DeleteCategoryByName(ctx context.Context, categoryName string, userID uint64) error
	DeleteCategoryList(ctx context.Context, userID uint64) error
}
