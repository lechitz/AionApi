package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/aion-api/internal/category/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SaveCategory persists a category in cache with TTL.
func (s *Store) SaveCategory(ctx context.Context, category domain.Category, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategorySave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "category"),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(category.ID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryIDKeyFormat, category.UserID, category.ID)

	// if no expiration passed, fall back to default
	if expiration <= 0 {
		expiration = CategoryExpirationDefault
	}

	// Convert to cache struct with JSON tags (musttag requirement)
	categoryJSON := fromDomainCategory(category)
	data, err := json.Marshal(categoryJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeCategory, commonkeys.CategoryID, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveCategoryToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategorySavedSuccessfully)
	return nil
}

// SaveCategoryByName persists a category in cache by name with TTL.
func (s *Store) SaveCategoryByName(ctx context.Context, category domain.Category, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategorySave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "category"),
		attribute.String(commonkeys.CategoryName, category.Name),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryNameKeyFormat, category.UserID, category.Name)

	// if no expiration passed, fall back to default
	if expiration <= 0 {
		expiration = CategoryExpirationDefault
	}

	// Convert to cache struct with JSON tags (musttag requirement)
	categoryJSON := fromDomainCategory(category)
	data, err := json.Marshal(categoryJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeCategory, commonkeys.CategoryName, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveCategoryToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategorySavedSuccessfully)
	return nil
}

// SaveCategoryList persists a category list in cache with TTL.
func (s *Store) SaveCategoryList(ctx context.Context, userID uint64, categories []domain.Category, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "category_list"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryListKeyFormat, userID)

	// if no expiration passed, fall back to default
	if expiration <= 0 {
		expiration = CategoryListExpirationDefault
	}

	// Convert to cache structs with JSON tags (musttag requirement)
	categoriesJSON := make([]categoryCache, len(categories))
	for i, cat := range categories {
		categoriesJSON[i] = fromDomainCategory(cat)
	}

	data, err := json.Marshal(categoriesJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeCategoryList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveCategoryListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, CategoryListSavedSuccessfully)
	return nil
}
