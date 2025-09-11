package repository

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// UpdateCategory updates a handler in the database based on its ID and user ID, updating only fields specified in the updateFields map.
func (c CategoryRepository) UpdateCategory(ctx context.Context, categoryID uint64, userID uint64, updateFields map[string]interface{}) (domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "Update", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String("operation", "update"),
	))
	defer span.End()

	delete(updateFields, commonkeys.CategoryCreatedAt)

	var categoryDB model.CategoryDB
	if err := c.db.WithContext(ctx).
		Model(&categoryDB).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(updateFields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, err
	}

	if err := c.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, err
	}

	span.SetStatus(codes.Ok, "handler updated successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}
