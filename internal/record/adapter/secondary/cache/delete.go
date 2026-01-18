package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteRecord removes a record from cache by ID.
func (s *Store) DeleteRecord(ctx context.Context, recordID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordDelete, trace.WithAttributes(
		attribute.String("record_id", strconv.FormatUint(recordID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "record"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordIDKeyFormat, userID, recordID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteRecordFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordDeletedSuccessfully)
	return nil
}

// DeleteRecordsByDay removes records for a specific day from cache.
func (s *Store) DeleteRecordsByDay(ctx context.Context, userID uint64, date time.Time) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "record_by_day"),
	))
	defer span.End()

	dateStr := date.Format("2006-01-02")
	cacheKey := fmt.Sprintf(RecordDayKeyFormat, userID, dateStr)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListDeletedSuccessfully)
	return nil
}

// DeleteRecordsByCategory removes records by category from cache.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) DeleteRecordsByCategory(ctx context.Context, categoryID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListDelete, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "record_by_category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByCategoryKeyFormat, categoryID, userID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListDeletedSuccessfully)
	return nil
}

// DeleteRecordsByTag removes records by tag from cache.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) DeleteRecordsByTag(ctx context.Context, tagID, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListDelete, trace.WithAttributes(
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "record_by_tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByTagKeyFormat, tagID, userID)

	if err := s.cache.Del(ctx, cacheKey); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListDeletedSuccessfully)
	return nil
}
