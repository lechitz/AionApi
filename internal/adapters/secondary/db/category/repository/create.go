package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/category/mapper"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create creates a new category in the database and returns the created category or an error if the operation fails.
func (c CategoryRepository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer("CategoryRepository")
	ctx, span := tr.Start(ctx, "Create", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(commonkeys.CategoryName, category.Name),
		attribute.String("operation", "create"),
	))
	defer span.End()

	categoryDB := mapper.CategoryToDB(category)
	if err := c.db.WithContext(ctx).
		Create(&categoryDB).Error; err != nil {
		wrappedErr := fmt.Errorf("error creating category: %w", err)
		span.SetStatus(codes.Error, wrappedErr.Error())
		span.RecordError(wrappedErr)
		c.logger.Errorw("error creating category", commonkeys.Category, category, commonkeys.Error, wrappedErr.Error())
		return domain.Category{}, wrappedErr
	}

	span.SetStatus(codes.Ok, "category created successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}
