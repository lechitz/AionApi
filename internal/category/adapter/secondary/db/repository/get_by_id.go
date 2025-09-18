package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GetByID loads a category by its ID and user ID. Returns not-found or generic errors using constants.
func (c CategoryRepository) GetByID(ctx context.Context, categoryID uint64, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByIDRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.Operation, OpGetByID),
	))
	defer span.End()

	var categoryDB model.CategoryDB
	if err := c.db.WithContext(ctx).
		Select("category_id, user_id, name, description, color_hex, icon, created_at, updated_at").
		Where("category_id = ? AND user_id = ?", categoryID, userID).
		First(&categoryDB).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Error, ErrCategoryNotFoundMsg)
			span.RecordError(errors.New(ErrCategoryNotFoundMsg))
			return domain.Category{}, errors.New(ErrCategoryNotFoundMsg)
		}
		span.SetStatus(codes.Error, ErrGetCategoryMsg)
		span.RecordError(err)
		return domain.Category{}, errors.New(ErrGetCategoryMsg)
	}

	span.SetStatus(codes.Ok, StatusRetrievedByID)
	return mapper.CategoryFromDB(categoryDB), nil
}
