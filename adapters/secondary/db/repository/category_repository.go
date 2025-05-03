package repository

import (
	"context"
	"fmt"
	"github.com/lechitz/AionApi/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/adapters/secondary/db/model"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

func NewCategoryRepository(db *gorm.DB, logger logger.Logger) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}

func (c CategoryRepository) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {

	categoryDB := mapper.CategoryToDB(category)

	if err := c.db.WithContext(ctx).
		Create(&categoryDB).Error; err != nil {
		wrappedErr := fmt.Errorf("error creating category: %w", err)
		c.logger.Errorw("error creating category", "category", category, "error", wrappedErr.Error())
		return domain.Category{}, wrappedErr
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

func (c CategoryRepository) GetCategoryByID(ctx context.Context, categoryID uint64) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ?", categoryID).
		First(&categoryDB).Error; err != nil {
		c.logger.Errorw("error getting category", "id", categoryID, "error", err.Error())
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

func (c CategoryRepository) GetCategoryByName(ctx context.Context, name string) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("name = ?", name).
		First(&categoryDB).Error; err != nil {
		c.logger.Errorw("error getting category", "name", name, "error", err.Error())
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

func (c CategoryRepository) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	var categoriesDB []model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Find(&categoriesDB).Error; err != nil {
		c.logger.Errorw("error getting categories", "error", err.Error())
		return nil, err
	}

	categories := make([]domain.Category, len(categoriesDB))
	for i, categoryDB := range categoriesDB {
		categories[i] = mapper.CategoryFromDB(categoryDB)
	}

	return categories, nil
}
