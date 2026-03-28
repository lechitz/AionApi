package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteCategory removes a category from cache by ID.
func (s *Store) DeleteCategory(ctx context.Context, categoryID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryDelete, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryIDKeyFormat, userID, categoryID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteCategoryFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategoryDeletedSuccessfully)
	return nil
}

// DeleteCategoryByName removes a category from cache by name.
func (s *Store) DeleteCategoryByName(ctx context.Context, categoryName string, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryDelete, trace.WithAttributes(
		attribute.String(commonkeys.CategoryName, categoryName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryNameKeyFormat, userID, categoryName)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteCategoryFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategoryDeletedSuccessfully)
	return nil
}

// DeleteCategoryList removes a category list from cache.
func (s *Store) DeleteCategoryList(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryListDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "category_list"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryListKeyFormat, userID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteCategoryListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategoryListDeletedSuccessfully)
	return nil
}
