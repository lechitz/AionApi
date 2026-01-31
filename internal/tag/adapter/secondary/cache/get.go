package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/tag/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// tagCache is an intermediate struct with JSON tags for safe unmarshaling (musttag requirement).
type tagCache struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Icon        string     `json:"icon"`
	CategoryID  uint64     `json:"category_id"`
	ID          uint64     `json:"id"`
	UserID      uint64     `json:"user_id"`
}

// toDomain converts tagCache to domain.Tag.
func (t tagCache) toDomain() domain.Tag {
	return domain.Tag{
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
		Name:        t.Name,
		Description: t.Description,
		Icon:        t.Icon,
		CategoryID:  t.CategoryID,
		ID:          t.ID,
		UserID:      t.UserID,
	}
}

// GetTag retrieves a tag from cache by ID.
func (s *Store) GetTag(ctx context.Context, tagID, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagGet, trace.WithAttributes(
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagIDKeyFormat, userID, tagID)

	tagValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTagFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Tag{}, err
	}

	if tagValue == "" {
		span.SetStatus(codes.Ok, "tag not found in cache")
		return domain.Tag{}, nil
	}

	var tagJSON tagCache
	if err := json.Unmarshal([]byte(tagValue), &tagJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeTag, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Tag{}, err
	}

	span.SetStatus(codes.Ok, TagRetrievedSuccessfully)
	return tagJSON.toDomain(), nil
}

// GetTagByName retrieves a tag from cache by name.
func (s *Store) GetTagByName(ctx context.Context, tagName string, userID uint64) (domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagGet, trace.WithAttributes(
		attribute.String(commonkeys.TagName, tagName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagNameKeyFormat, userID, tagName)

	tagValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTagFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Tag{}, err
	}

	if tagValue == "" {
		span.SetStatus(codes.Ok, "tag not found in cache")
		return domain.Tag{}, nil
	}

	var tagJSON tagCache
	if err := json.Unmarshal([]byte(tagValue), &tagJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeTag, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Tag{}, err
	}

	span.SetStatus(codes.Ok, TagRetrievedSuccessfully)
	return tagJSON.toDomain(), nil
}

// GetTagList retrieves a tag list from cache.
func (s *Store) GetTagList(ctx context.Context, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "tag_list"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagListKeyFormat, userID)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTagListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "tag list not found in cache")
		return nil, nil
	}

	var tagsJSON []tagCache
	if err := json.Unmarshal([]byte(listValue), &tagsJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeTagList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	// Convert to domain slice
	tags := make([]domain.Tag, len(tagsJSON))
	for i, tagJSON := range tagsJSON {
		tags[i] = tagJSON.toDomain()
	}

	span.SetStatus(codes.Ok, TagListRetrievedSuccessfully)
	return tags, nil
}

// GetTagsByCategory retrieves tags by category from cache.
func (s *Store) GetTagsByCategory(ctx context.Context, categoryID, userID uint64) ([]domain.Tag, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListGet, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "tag_by_category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagByCategoryKeyFormat, categoryID, userID)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTagListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "tag list by category not found in cache")
		return nil, nil
	}

	var tagsJSON []tagCache
	if err := json.Unmarshal([]byte(listValue), &tagsJSON); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeTagList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	// Convert to domain slice
	tags := make([]domain.Tag, len(tagsJSON))
	for i, tagJSON := range tagsJSON {
		tags[i] = tagJSON.toDomain()
	}

	span.SetStatus(codes.Ok, TagListRetrievedSuccessfully)
	return tags, nil
}
