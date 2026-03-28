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

// categoryCache is an intermediate struct with JSON tags for safe unmarshaling (musttag requirement).
type categoryCache struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Color       string     `json:"color"`
	Icon        string     `json:"icon"`
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
}

// toDomain converts categoryCache to domain.Category.
func (c categoryCache) toDomain() domain.Category {
	return domain.Category{
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
		Name:        c.Name,
		Description: c.Description,
		Color:       c.Color,
		Icon:        c.Icon,
		ID:          c.ID,
		UserID:      c.UserID,
	}
}

// fromDomainCategory converts domain.Category to categoryCache (for Marshal).
func fromDomainCategory(c domain.Category) categoryCache {
	return categoryCache{
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		DeletedAt:   c.DeletedAt,
		Name:        c.Name,
		Description: c.Description,
		Color:       c.Color,
		Icon:        c.Icon,
		ID:          c.ID,
		UserID:      c.UserID,
	}
}

// GetCategory retrieves a category from cache by ID.
func (s *Store) GetCategory(ctx context.Context, categoryID, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryGet, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryIDKeyFormat, userID, categoryID)

	categoryValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetCategoryFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Category{}, err
	}

	if categoryValue == "" {
		span.SetStatus(codes.Ok, "category not found in cache")
		return domain.Category{}, nil
	}

	var categoryJSON categoryCache
	if err := json.Unmarshal([]byte(categoryValue), &categoryJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeCategory, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.SetStatus(codes.Ok, CategoryRetrievedSuccessfully)
	return categoryJSON.toDomain(), nil
}

// GetCategoryByName retrieves a category from cache by name.
func (s *Store) GetCategoryByName(ctx context.Context, categoryName string, userID uint64) (domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryGet, trace.WithAttributes(
		attribute.String(commonkeys.CategoryName, categoryName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryNameKeyFormat, userID, categoryName)

	categoryValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetCategoryFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Category{}, err
	}

	if categoryValue == "" {
		span.SetStatus(codes.Ok, "category not found in cache")
		return domain.Category{}, nil
	}

	var categoryJSON categoryCache
	if err := json.Unmarshal([]byte(categoryValue), &categoryJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeCategory, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Category{}, err
	}

	span.SetStatus(codes.Ok, CategoryRetrievedSuccessfully)
	return categoryJSON.toDomain(), nil
}

// GetCategoryList retrieves a category list from cache.
func (s *Store) GetCategoryList(ctx context.Context, userID uint64) ([]domain.Category, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameCategoryListGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "category_list"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(CategoryListKeyFormat, userID)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetCategoryListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "category list not found in cache")
		return nil, nil
	}

	var categoriesJSON []categoryCache
	if err := json.Unmarshal([]byte(listValue), &categoriesJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeCategoryList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	// Convert to domain slice
	categories := make([]domain.Category, len(categoriesJSON))
	for i, catJSON := range categoriesJSON {
		categories[i] = catJSON.toDomain()
	}

	span.SetStatus(codes.Ok, CategoryListRetrievedSuccessfully)
	return categories, nil
}
