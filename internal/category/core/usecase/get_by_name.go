package usecase

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByName retrieves a handler by its name from the database and returns it.
func (s *Service) GetByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetCategoryByName)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetCategoryByName),
		attribute.String(commonkeys.CategoryName, categoryName),
	)

	if categoryName == "" {
		span.SetStatus(codes.Error, CategoryNameIsRequired)
		s.Logger.ErrorwCtx(ctx, CategoryNameIsRequired, commonkeys.CategoryName, categoryName)
		return domain.Category{}, errors.New(CategoryNameIsRequired)
	}

	span.AddEvent(EventRepositoryGet)
	category, err := s.CategoryRepository.GetByName(ctx, categoryName, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetCategoryByName)
		s.Logger.ErrorwCtx(ctx, FailedToGetCategoryByName, commonkeys.CategoryName, categoryName, commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrievedByName)
	return category, nil
}
