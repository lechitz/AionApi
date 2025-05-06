package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/adapters/secondary/db/constants"
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

func (c CategoryRepository) GetCategoryByID(ctx context.Context, category domain.Category) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		First(&categoryDB).Error; err != nil {
		c.logger.Errorw("error getting category", "category_id", category.ID, "error", err.Error())

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, fmt.Errorf("category not found")
		}
		return domain.Category{}, fmt.Errorf("error getting category")
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

func (c CategoryRepository) GetCategoryByName(ctx context.Context, category domain.Category) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("name = ? AND user_id = ?", category.Name, category.UserID).
		First(&categoryDB).Error; err != nil {
		c.logger.Errorw("error getting category", "name", category.Name, "error", err.Error())
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

func (c CategoryRepository) GetAllCategories(ctx context.Context, userID uint64) ([]domain.Category, error) {
	var categoriesDB []model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("user_id = ?", userID).
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

func (c CategoryRepository) UpdateCategory(ctx context.Context, categoryID uint64, userID uint64, updateFields map[string]interface{}) (domain.Category, error) {
	delete(updateFields, constants.CreatedAt)

	var categoryDB model.CategoryDB
	if err := c.db.WithContext(ctx).
		Model(&categoryDB).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(updateFields).Error; err != nil {
		c.logger.Errorw("error updating category", "category_id", categoryID, "user_id", userID, "error", err.Error())
		return domain.Category{}, err
	}

	if err := c.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error; err != nil {
		c.logger.Errorw("error fetching updated category", "category_id", categoryID, "user_id", userID, "error", err.Error())
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}
