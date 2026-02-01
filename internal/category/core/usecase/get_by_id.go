package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/category/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID retrieves a handler by its ID from the database and returns it.
func (s *Service) GetByID(ctx context.Context, categoryID, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetCategoryByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetCategoryByID),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
	)

	if categoryID == 0 {
		span.SetStatus(codes.Error, CategoryIDIsRequired)
		s.Logger.ErrorwCtx(ctx, CategoryIDIsRequired, commonkeys.CategoryID, strconv.FormatUint(categoryID, 10))
		return domain.Category{}, errors.New(CategoryIDIsRequired)
	}

	span.AddEvent(EventCheckCache)
	cachedCategory, err := s.CategoryCache.GetCategory(ctx, categoryID, userID)
	if err == nil && cachedCategory.ID != 0 {
		span.AddEvent(EventCacheHit)
		span.SetStatus(codes.Ok, StatusRetrievedByID)
		return cachedCategory, nil
	}

	span.AddEvent(EventRepositoryGet)
	categoryDB, err := s.CategoryRepository.GetByID(ctx, categoryID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetCategoryByID)
		s.Logger.ErrorwCtx(ctx, FailedToGetCategoryByID, commonkeys.CategoryID, strconv.FormatUint(categoryID, 10), commonkeys.Error, err.Error())
		return domain.Category{}, errors.New(FailedToGetCategoryByID)
	}

	span.AddEvent(EventSaveToCache)
	err = s.CategoryCache.SaveCategory(ctx, categoryDB, 0) // use default TTL
	if err != nil {
		s.Logger.WarnwCtx(ctx, WarnFailedToSaveCategoryToCache,
			commonkeys.CategoryID, categoryDB.ID,
			commonkeys.Error, err.Error(),
		)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrievedByID)
	return categoryDB, nil
}
