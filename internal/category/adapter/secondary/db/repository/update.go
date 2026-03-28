package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/aion-api/internal/category/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// UpdateCategory updates a handler in the database based on its ID and user ID, updating only fields specified in the updateFields map.
func (c CategoryRepository) UpdateCategory(ctx context.Context, categoryID uint64, userID uint64, updateFields map[string]interface{}) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdateRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.Operation, OpUpdate),
	))
	defer span.End()

	delete(updateFields, commonkeys.CategoryCreatedAt)
	if len(updateFields) == 0 {
		span.SetStatus(codes.Error, "no fields to update")
		return domain.Category{}, errors.New("no fields to update")
	}

	var categoryDB model.CategoryDB
	q := c.db.WithContext(ctx).
		Model(&categoryDB).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		Updates(updateFields)
	if err := q.Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, err
	}

	// If no rows were affected, the category either doesn't exist or doesn't belong to the user.
	if q.RowsAffected() == 0 {
		err := gorm.ErrRecordNotFound
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, errors.New(ErrCategoryNotFoundMsg)
	}

	if err := c.db.WithContext(ctx).
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, err
	}

	span.SetStatus(codes.Ok, StatusUpdated)
	return mapper.CategoryFromDB(categoryDB), nil
}
