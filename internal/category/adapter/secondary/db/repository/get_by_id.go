package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/model"
	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GetByID retrieves a handler by its ID and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetByID(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetByID", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
		attribute.String("operation", "get_by_id"),
	))
	defer span.End()

	var categoryDB model.CategoryDB

	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ? AND user_id = ?", category.ID, category.UserID).
		First(&categoryDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, "handler not found")
			span.RecordError(errors.New("handler not found"))
			return domain.Category{}, errors.New("handler not found")
		}
		span.SetStatus(codes.Error, "error getting handler")
		span.RecordError(err)
		return domain.Category{}, errors.New("error getting handler")
	}

	span.SetStatus(codes.Ok, "handler retrieved by id successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}
