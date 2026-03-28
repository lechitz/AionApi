package controller

import (
	"context"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/lechitz/aion-api/internal/adapter/primary/graphql/model"
	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	maxInt32Value = int64(1<<31 - 1)
	minInt32Value = int64(-1 << 31)
)

// RecordStats computes aggregated statistics over records using the same filter semantics from search.
func (c *controller) RecordStats(ctx context.Context, filters *model.RecordStatsFilters, userID uint64) (*model.RecordStats, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanRecordStats)
	defer span.End()

	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	if userID == 0 {
		span.SetStatus(codes.Error, ErrUserIDNotFound.Error())
		return nil, ErrUserIDNotFound
	}

	domainFilters := convertRecordStatsFilters(filters)
	records, err := c.RecordService.SearchRecords(ctx, userID, domainFilters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, MsgStatsError)
		c.Logger.ErrorwCtx(ctx, MsgStatsError,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(userID, 10),
		)
		return nil, err
	}

	stats := buildRecordStats(records)
	span.SetAttributes(attribute.Int(AttrRecordsCount, int(stats.TotalRecords)))
	span.SetStatus(codes.Ok, StatusStatsComputed)
	return stats, nil
}

func convertRecordStatsFilters(filters *model.RecordStatsFilters) domain.SearchFilters {
	out := domain.SearchFilters{
		Query:  "",
		Limit:  5000,
		Offset: 0,
	}

	if filters == nil {
		return out
	}

	if filters.Query != nil {
		out.Query = strings.TrimSpace(*filters.Query)
	}

	out.CategoryIDs = convertIDSlice(filters.CategoryIds)
	out.TagIDs = convertIDSlice(filters.TagIds)

	if filters.StartDate != nil {
		if t, err := time.Parse(time.RFC3339, *filters.StartDate); err == nil {
			out.StartDate = &t
		}
	}
	if filters.EndDate != nil {
		if t, err := time.Parse(time.RFC3339, *filters.EndDate); err == nil {
			out.EndDate = &t
		}
	}

	if filters.Limit != nil {
		out.Limit = int(*filters.Limit)
	}

	if out.Limit <= 0 {
		out.Limit = 5000
	}
	if out.Limit > 20000 {
		out.Limit = 20000
	}

	return out
}

func buildRecordStats(records []domain.Record) *model.RecordStats {
	total := len(records)
	if total == 0 {
		return &model.RecordStats{
			TotalRecords:         0,
			RecordsWithValue:     0,
			TotalDurationSeconds: 0,
			SumValue:             0,
			AvgValue:             0,
			AvgDurationSeconds:   0,
			MinValue:             nil,
			MaxValue:             nil,
		}
	}

	var valueCount int
	var sumValue float64
	var sumDuration int
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64

	for _, record := range records {
		if record.Value != nil {
			v := *record.Value
			sumValue += v
			valueCount++
			if v < minValue {
				minValue = v
			}
			if v > maxValue {
				maxValue = v
			}
		}
		if record.DurationSecs != nil {
			sumDuration += *record.DurationSecs
		}
	}

	avgDuration := float64(sumDuration) / float64(total)
	avgValue := 0.0
	if valueCount > 0 {
		avgValue = sumValue / float64(valueCount)
	}

	var minOut *float64
	var maxOut *float64
	if valueCount > 0 {
		minOut = &minValue
		maxOut = &maxValue
	}

	return &model.RecordStats{
		TotalRecords:         safeIntToInt32(total),
		RecordsWithValue:     safeIntToInt32(valueCount),
		TotalDurationSeconds: safeIntToInt32(sumDuration),
		SumValue:             sumValue,
		AvgValue:             avgValue,
		AvgDurationSeconds:   avgDuration,
		MinValue:             minOut,
		MaxValue:             maxOut,
	}
}

func safeIntToInt32(value int) int32 {
	v64 := int64(value)
	if v64 > maxInt32Value {
		return int32(maxInt32Value)
	}
	if v64 < minInt32Value {
		return int32(minInt32Value)
	}
	return int32(v64)
}
