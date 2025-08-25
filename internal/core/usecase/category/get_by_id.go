package category

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/category/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID retrieves a category by its ID from the database and returns it.
func (s *Service) GetByID(ctx context.Context, categoryID, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(constants.TracerName)
	ctx, span := tr.Start(ctx, constants.SpanGetCategoryByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanGetCategoryByID),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	if categoryID == 0 {
		span.SetStatus(codes.Error, constants.CategoryIDIsRequired)
		s.Logger.ErrorwCtx(ctx, constants.CategoryIDIsRequired, commonkeys.CategoryID, strconv.FormatUint(categoryID, 10))
		return domain.Category{}, errors.New(constants.CategoryIDIsRequired)
	}

	span.AddEvent(constants.EventRepositoryGet)
	categoryDB, err := s.CategoryRepository.GetByID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.FailedToGetCategoryByID)
		s.Logger.ErrorwCtx(ctx, constants.FailedToGetCategoryByID, commonkeys.CategoryID, strconv.FormatUint(categoryID, 10), commonkeys.Error, err.Error())
		return domain.Category{}, errors.New(constants.FailedToGetCategoryByID)
	}

	span.AddEvent(constants.EventSuccess)
	span.SetStatus(codes.Ok, "retrieved")
	return categoryDB, nil
}
