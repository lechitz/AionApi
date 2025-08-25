package category

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByName retrieves a category by its name from the database and returns it.
func (s *Service) GetByName(ctx context.Context, categoryName string) (domain.Category, error) {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanGetCategoryByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetCategoryByName),
		attribute.String(commonkeys.CategoryName, categoryName),
	)

	if categoryName == "" {
		span.SetStatus(codes.Error, constants.CategoryNameIsRequired)
		s.Logger.ErrorwCtx(ctx, constants.CategoryNameIsRequired, commonkeys.CategoryName, categoryName)
		return domain.Category{}, errors.New(constants.CategoryNameIsRequired)
	}

	span.AddEvent(constants.EventRepositoryGet)
	category, err := s.CategoryRepository.GetByName(ctx, categoryName)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToGetCategoryByName)
		s.Logger.ErrorwCtx(ctx, constants.FailedToGetCategoryByName, commonkeys.CategoryName, categoryName, commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "retrieved")
	return category, nil
}
