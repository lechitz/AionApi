package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/category/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create creates a new handler in the database and returns the created handler or an error if the operation fails.
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
		wrappedErr := fmt.Errorf("error creating handler: %w", err)
		span.SetStatus(codes.Error, wrappedErr.Error())
		span.RecordError(wrappedErr)
		c.logger.Errorw("error creating handler", commonkeys.Category, category, commonkeys.Error, wrappedErr.Error())
		return domain.Category{}, wrappedErr
	}

	span.SetStatus(codes.Ok, "handler created successfully")
	return mapper.CategoryFromDB(categoryDB), nil
}
