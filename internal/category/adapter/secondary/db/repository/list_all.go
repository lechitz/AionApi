package repository

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/model"
	"github.com/lechitz/AionApi/internal/core/category/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ListAll retrieves all categories associated with a specific user defined by the userID. Returns a slice of domain.Category or an error.
func (c CategoryRepository) ListAll(ctx context.Context, userID uint64) ([]domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "ListAll", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
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

	categories := make([]domain.Category, len(categoriesDB))
	for i, categoryDB := range categoriesDB {
		categories[i] = mapper.CategoryFromDB(categoryDB)
	}

	span.SetStatus(codes.Ok, "all categories retrieved successfully")
	return categories, nil
}
