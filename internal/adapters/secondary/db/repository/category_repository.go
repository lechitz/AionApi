// Package repository provides methods for interacting with the category database.
package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"gorm.io/gorm"
)

// CategoryRepository manages database operations related to category entities.
// It uses gorm.DB for ORM and logger.Logger for logging operations.
type CategoryRepository struct {
	db     *gorm.DB
	logger logger.Logger
}

// NewCategoryRepository creates a new instance of CategoryRepository with a given gorm.DB and logger.
func NewCategoryRepository(db *gorm.DB, logger logger.Logger) *CategoryRepository {
	return &CategoryRepository{
		db:     db,
		logger: logger,
	}
}

// CreateCategory creates a new category in the database and returns the created category or an error if the operation fails.
func (c CategoryRepository) CreateCategory(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	categoryDB := mapper.CategoryToDB(category)

	if err := c.db.WithContext(ctx).
		Create(&categoryDB).Error; err != nil {
		wrappedErr := fmt.Errorf("error creating category: %w", err)
		c.logger.Errorw(
			"error creating category",
			"category",
			category,
			"error",
			wrappedErr.Error(),
		)
		return domain.Category{}, wrappedErr
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

// GetCategoryByID retrieves a category by its ID and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetCategoryByID(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		First(&categoryDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Category{}, errors.New("category not found")
		}
		return domain.Category{}, errors.New("error getting category")
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

// GetCategoryByName retrieves a category by its name and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetCategoryByName(
	ctx context.Context,
	category domain.Category,
) (domain.Category, error) {
	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("name = ? AND user_id = ?", category.Name, category.UserID).
		First(&categoryDB).Error; err != nil {
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

// GetAllCategories retrieves all categories associated with a specific user defined by the userID. Returns a slice of domain.Category or an error.
func (c CategoryRepository) GetAllCategories(
	ctx context.Context,
	userID uint64,
) ([]domain.Category, error) {
	var categoriesDB []model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("user_id = ?", userID).
		Find(&categoriesDB).Error; err != nil {
		return nil, err
	}

	categories := make([]domain.Category, len(categoriesDB))
	for i, categoryDB := range categoriesDB {
		categories[i] = mapper.CategoryFromDB(categoryDB)
	}

	return categories, nil
}

// UpdateCategory updates a category in the database based on its ID and user ID, updating only fields specified in the updateFields map.
func (c CategoryRepository) UpdateCategory(
	ctx context.Context, categoryID uint64,
	userID uint64,
	updateFields map[string]interface{},
) (domain.Category, error) {
	delete(updateFields, constants.CreatedAt)

	var categoryDB model.CategoryDB
	if err := c.db.WithContext(ctx).
		Model(&categoryDB).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(updateFields).Error; err != nil {
		return domain.Category{}, err
	}

	if err := c.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error; err != nil {
		return domain.Category{}, err
	}

	return mapper.CategoryFromDB(categoryDB), nil
}

// SoftDeleteCategory updates the DeletedAt and UpdatedAt fields to mark a category as soft-deleted based on category ID and user ID.
func (c CategoryRepository) SoftDeleteCategory(
	ctx context.Context,
	category domain.Category,
) error {
	fields := map[string]interface{}{
		constants.DeletedAt: time.Now().UTC(),
		constants.UpdatedAt: time.Now().UTC(),
	}

	if err := c.db.WithContext(ctx).
		Model(&model.CategoryDB{}).
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		Updates(fields).Error; err != nil {
		return err
	}

	return nil
}
