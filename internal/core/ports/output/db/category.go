package db

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type CategoryCreator interface {
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

type CategoryRetriver interface {
	GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error)
	GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error)
	GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error)
}

type CategoryUpdater interface {
	UpdateCategory(ctx context.Context, categoryID uint64, userID uint64, fields map[string]interface{}) (domain.Category, error)
}

type CategoryDeleter interface {
	SoftDeleteCategory(ctx context.Context, category domain.Category) error
}

type CategoryStore interface {
	CategoryCreator
	CategoryRetriver
	CategoryUpdater
	CategoryDeleter
}
