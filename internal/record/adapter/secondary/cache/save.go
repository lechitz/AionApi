package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SaveRecord persists a record in cache with TTL.
func (s *Store) SaveRecord(ctx context.Context, record domain.Record, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "record"),
		attribute.String("record_id", strconv.FormatUint(record.ID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordIDKeyFormat, record.UserID, record.ID)

	if expiration <= 0 {
		expiration = RecordExpirationDefault
	}

	data, err := json.Marshal(record)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeRecord, "record_id", cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveRecordToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordSavedSuccessfully)
	return nil
}

// SaveRecordsByDay persists records for a specific day in cache with TTL.
func (s *Store) SaveRecordsByDay(ctx context.Context, userID uint64, date time.Time, records []domain.Record, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "record_by_day"),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	dateStr := date.Format("2006-01-02")
	cacheKey := fmt.Sprintf(RecordDayKeyFormat, userID, dateStr)

	if expiration <= 0 {
		expiration = RecordListExpirationDefault
	}

	data, err := json.Marshal(records)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveRecordListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListSavedSuccessfully)
	return nil
}

// SaveRecordsByCategory persists records by category in cache with TTL.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) SaveRecordsByCategory(ctx context.Context, categoryID, userID uint64, records []domain.Record, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "record_by_category"),
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByCategoryKeyFormat, categoryID, userID)

	if expiration <= 0 {
		expiration = RecordListExpirationDefault
	}

	data, err := json.Marshal(records)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveRecordListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListSavedSuccessfully)
	return nil
}

// SaveRecordsByTag persists records by tag in cache with TTL.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) SaveRecordsByTag(ctx context.Context, tagID, userID uint64, records []domain.Record, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, "record_by_tag"),
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByTagKeyFormat, tagID, userID)

	if expiration <= 0 {
		expiration = RecordListExpirationDefault
	}

	data, err := json.Marshal(records)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	if err := s.cache.Set(ctx, cacheKey, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveRecordListToCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, RecordListSavedSuccessfully)
	return nil
}
