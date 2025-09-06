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

// GetByName retrieves a handler by its name and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetByName(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetByName", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(commonkeys.CategoryName, category.Name),
		attribute.String("operation", "get_by_name"),
	))
	defer span.End()

	var categoryDB model.CategoryDB
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", category.UserID, category.Name).
		First(&categoryDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "handler not found (normal case)")
			return domain.Category{}, nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.Category{}, err
	}

	span.SetStatus(codes.Ok, "handler fetched successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}
