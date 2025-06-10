package category

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
)

type CategoryService struct {
	CategoryRepository db.CategoryStore
	Logger             logger.Logger
}

func NewCategoryService(
	categoryRepository db.CategoryStore,
	logger logger.Logger,
) *CategoryService {
	return &CategoryService{
		CategoryRepository: categoryRepository,
		Logger:             logger,
	}
}
