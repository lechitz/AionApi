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

// GetByName retrieves a handler by its name and user ID from the database and returns it as a domain.Category or an error if not found.
func (c CategoryRepository) GetByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "GetByName", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.CategoryName, categoryName),
		attribute.String("operation", "get_by_name"),
	))
	defer span.End()

	var categoryDB model.CategoryDB
	err := c.db.WithContext(ctx).
		Where("user_id = ? AND name = ?", userID, categoryName).
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
