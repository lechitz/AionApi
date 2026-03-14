package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetRecord retrieves a record from cache by ID.
func (s *Store) GetRecord(ctx context.Context, recordID, userID uint64) (domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordGet, trace.WithAttributes(
		attribute.String("record_id", strconv.FormatUint(recordID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "record"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordIDKeyFormat, userID, recordID)

	recordValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetRecordFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Record{}, err
	}

	if recordValue == "" {
		span.SetStatus(codes.Ok, "record not found in cache")
		return domain.Record{}, nil
	}

	var record domain.Record
	if err := json.Unmarshal([]byte(recordValue), &record); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeRecord, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.Record{}, err
	}

	span.SetStatus(codes.Ok, RecordRetrievedSuccessfully)
	return record, nil
}

// GetRecordsByDay retrieves records for a specific day from cache.
func (s *Store) GetRecordsByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "record_by_day"),
	))
	defer span.End()

	dateStr := date.Format("2006-01-02")
	cacheKey := fmt.Sprintf(RecordDayKeyFormat, userID, dateStr)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "record list by day not found in cache")
		return nil, nil
	}

	var records []domain.Record
	if err := json.Unmarshal([]byte(listValue), &records); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	span.SetStatus(codes.Ok, RecordListRetrievedSuccessfully)
	return records, nil
}

// GetRecordsByCategory retrieves records by category from cache.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) GetRecordsByCategory(ctx context.Context, categoryID, userID uint64) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListGet, trace.WithAttributes(
		attribute.String(commonkeys.CategoryID, strconv.FormatUint(categoryID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "record_by_category"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByCategoryKeyFormat, categoryID, userID)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "record list by category not found in cache")
		return nil, nil
	}

	var records []domain.Record
	if err := json.Unmarshal([]byte(listValue), &records); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	span.SetStatus(codes.Ok, RecordListRetrievedSuccessfully)
	return records, nil
}

// GetRecordsByTag retrieves records by tag from cache.
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) GetRecordsByTag(ctx context.Context, tagID, userID uint64) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameRecordListGet, trace.WithAttributes(
		attribute.String(commonkeys.TagID, strconv.FormatUint(tagID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "record_by_tag"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(RecordByTagKeyFormat, tagID, userID)

	listValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetRecordListFromCache, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	if listValue == "" {
		span.SetStatus(codes.Ok, "record list by tag not found in cache")
		return nil, nil
	}

	var records []domain.Record
	if err := json.Unmarshal([]byte(listValue), &records); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeRecordList, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return nil, err
	}

	span.SetStatus(codes.Ok, RecordListRetrievedSuccessfully)
	return records, nil
}
