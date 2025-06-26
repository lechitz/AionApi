// Package repository provides methods for interacting with the category database.
package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/constants"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"gorm.io/gorm"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
func (c CategoryRepository) CreateCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "CreateCategory", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(constants.CategoryName, category.Name),
		attribute.String("operation", "create"),
	))
	defer span.End()

	categoryDB := mapper.CategoryToDB(category)
	if err := c.db.WithContext(ctx).
		Create(&categoryDB).Error; err != nil {
		wrappedErr := fmt.Errorf("error creating category: %w", err)
		span.SetStatus(codes.Error, wrappedErr.Error())
		span.RecordError(wrappedErr)
		c.logger.Errorw(
			"error creating category",
			"category",
			category,
			"error",
			wrappedErr.Error(),
		)
		return entity.Category{}, wrappedErr
	}

	span.SetStatus(codes.Ok, "category created successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}

// GetCategoryByID retrieves a category by its ID and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetCategoryByID(ctx context.Context, category entity.Category) (entity.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetCategoryByID", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(constants.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String("operation", "get_by_id"),
	))
	defer span.End()

	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		First(&categoryDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, "category not found")
			span.RecordError(errors.New("category not found"))
			return entity.Category{}, errors.New("category not found")
		}
		span.SetStatus(codes.Error, "error getting category")
		span.RecordError(err)
		return entity.Category{}, errors.New("error getting category")
	}

	span.SetStatus(codes.Ok, "category retrieved by id successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}

// GetCategoryByName retrieves a category by its name and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetCategoryByName(ctx context.Context, category entity.Category) (entity.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetCategoryByName", trace.WithAttributes(
		attribute.String("user_id", strconv.FormatUint(category.UserID, 10)),
		attribute.String("category_name", category.Name),
		attribute.String("operation", "get_by_name"),
	))
	defer span.End()

	var categoryDB model.CategoryDB
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", category.UserID, category.Name).
		First(&categoryDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "category not found (normal case)")
			return entity.Category{}, nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return entity.Category{}, err
	}

	span.SetStatus(codes.Ok, "category fetched successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}

// GetAllCategories retrieves all categories associated with a specific user defined by the userID. Returns a slice of domain.Category or an error.
func (c CategoryRepository) GetAllCategories(ctx context.Context, userID uint64) ([]entity.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetAllCategories", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("operation", "get_all"),
	))
	defer span.End()

	var categoriesDB []model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("user_id = ?", userID).
		Find(&categoriesDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return nil, err
	}

	categories := make([]entity.Category, len(categoriesDB))
	for i, categoryDB := range categoriesDB {
		categories[i] = mapper.CategoryFromDB(categoryDB)
	}

	span.SetStatus(codes.Ok, "all categories retrieved successfully")
	return categories, nil
}

// UpdateCategory updates a category in the database based on its ID and user ID, updating only fields specified in the updateFields map.
func (c CategoryRepository) UpdateCategory(
	ctx context.Context,
	categoryID uint64,
	userID uint64,
	updateFields map[string]interface{},
) (entity.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "UpdateCategory", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(constants.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String("operation", "update"),
	))
	defer span.End()

	delete(updateFields, constants.CreatedAt)

	var categoryDB model.CategoryDB
	if err := c.db.WithContext(ctx).
		Model(&categoryDB).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(updateFields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return entity.Category{}, err
	}

	if err := c.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return entity.Category{}, err
	}

	span.SetStatus(codes.Ok, "category updated successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}

// SoftDeleteCategory updates the DeletedAt and UpdatedAt fields to mark a category as soft-deleted based on category ID and user ID.
func (c CategoryRepository) SoftDeleteCategory(
	ctx context.Context,
	category entity.Category,
) error {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "SoftDeleteCategory", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(constants.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String("operation", "soft_delete"),
	))
	defer span.End()

	fields := map[string]interface{}{
		constants.DeletedAt: time.Now().UTC(),
		constants.UpdatedAt: time.Now().UTC(),
	}

	if err := c.db.WithContext(ctx).
		Model(&model.CategoryDB{}).
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return err
	}

	span.SetStatus(codes.Ok, "category soft deleted successfully")
	return nil
}
