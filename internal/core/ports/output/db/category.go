package db

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
)

type CategoryCreator interface {
	CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error)
}

type CategoryRetriver interface {
	GetCategoryByID(ctx context.Context, id uint64) (domain.Category, error)
	GetCategoryByName(ctx context.Context, name string) (domain.Category, error)
	GetAllCategories(ctx context.Context) ([]domain.Category, error)
}

type CategoryStore interface {
	CategoryCreator
	CategoryRetriver
}
