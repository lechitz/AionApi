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

// Create inserts a new category and returns it, or an error on failure.
func (c CategoryRepository) Create(ctx context.Context, category domain.Category) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreateRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(category.UserID, 10)),
		attribute.String(commonkeys.CategoryName, category.Name),
		attribute.String(commonkeys.Operation, OpCreate),
	))
	defer span.End()

	categoryDB := mapper.CategoryToDB(category)
	if err := c.db.WithContext(ctx).Create(&categoryDB).Error; err != nil {
		wrappedErr := fmt.Errorf("%s: %w", ErrCreateCategoryMsg, err)
		span.SetStatus(codes.Error, wrappedErr.Error())
		span.RecordError(wrappedErr)
		c.logger.Errorw(ErrCreateCategoryMsg, commonkeys.Category, category, commonkeys.Error, wrappedErr.Error())
		return domain.Category{}, wrappedErr
	}

	span.SetStatus(codes.Ok, StatusCategoryCreated)
	return mapper.CategoryFromDB(categoryDB), nil
}
