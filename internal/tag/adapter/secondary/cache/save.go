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

// fromDomainTag converts domain.Tag to tagCache (for Marshal).
func fromDomainTag(t domain.Tag) tagCache {
	return tagCache{
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
		Name:        t.Name,
		Description: t.Description,
		CategoryID:  t.CategoryID,
		ID:          t.ID,
		UserID:      t.UserID,
	}
}

// SaveTag persists a tag in cache with TTL.
func (s *Store) SaveTag(ctx context.Context, tag domain.Tag, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "tag"),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tag.ID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagIDKeyFormat, tag.UserID, tag.ID)

	if expiration <= 0 {
		expiration = TagExpirationDefault
	}

	// Convert to cache struct with JSON tags (musttag requirement)
	tagJSON := fromDomainTag(tag)
	data, err := json.Marshal(tagJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeTag, commonkeys.TagID, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTagToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagSavedSuccessfully)
	return nil
}

// SaveTagByName persists a tag in cache by name with TTL.
func (s *Store) SaveTagByName(ctx context.Context, tag domain.Tag, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "tag"),
		attribute.String(commonkeys.TagName, tag.Name),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagNameKeyFormat, tag.UserID, tag.Name)

	if expiration <= 0 {
		expiration = TagExpirationDefault
	}

	// Convert to cache struct with JSON tags (musttag requirement)
	tagJSON := fromDomainTag(tag)
	data, err := json.Marshal(tagJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeTag, commonkeys.TagName, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTagToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagSavedSuccessfully)
	return nil
}

// SaveTagList persists a tag list in cache with TTL.
func (s *Store) SaveTagList(ctx context.Context, userID uint64, tags []domain.Tag, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "tag_list"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagListKeyFormat, userID)

	if expiration <= 0 {
		expiration = TagListExpirationDefault
	}

	// Convert to cache structs with JSON tags (musttag requirement)
	tagsJSON := make([]tagCache, len(tags))
	for i, tag := range tags {
		tagsJSON[i] = fromDomainTag(tag)
	}

	data, err := json.Marshal(tagsJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeTagList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTagListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagListSavedSuccessfully)
	return nil
}

// SaveTagsByCategory persists tags by category in cache with TTL.
func (s *Store) SaveTagsByCategory(ctx context.Context, categoryID, userID uint64, tags []domain.Tag, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "tag_by_category"),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagByCategoryKeyFormat, categoryID, userID)

	if expiration <= 0 {
		expiration = TagListExpirationDefault
	}

	// Convert to cache structs with JSON tags (musttag requirement)
	tagsJSON := make([]tagCache, len(tags))
	for i, tag := range tags {
		tagsJSON[i] = fromDomainTag(tag)
	}

	data, err := json.Marshal(tagsJSON)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeTagList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveTagListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagListSavedSuccessfully)
	return nil
}
