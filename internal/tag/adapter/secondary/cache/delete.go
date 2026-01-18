package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteTag removes a tag from cache by ID.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) DeleteTag(ctx context.Context, tagID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagDelete, trace.WithAttributes(
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagIDKeyFormat, userID, tagID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteTagFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagDeletedSuccessfully)
	return nil
}

// DeleteTagByName removes a tag from cache by name.
func (s *Store) DeleteTagByName(ctx context.Context, tagName string, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagDelete, trace.WithAttributes(
		attribute.String(commonkeys.TagName, tagName),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagNameKeyFormat, userID, tagName)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteTagFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagDeletedSuccessfully)
	return nil
}

// DeleteTagList removes a tag list from cache.
func (s *Store) DeleteTagList(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "tag_list"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagListKeyFormat, userID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteTagListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagListDeletedSuccessfully)
	return nil
}

// DeleteTagsByCategory removes tags by category from cache.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) DeleteTagsByCategory(ctx context.Context, categoryID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameTagListDelete, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "tag_by_category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TagByCategoryKeyFormat, categoryID, userID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteTagListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, TagListDeletedSuccessfully)
	return nil
}
