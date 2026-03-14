package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListByDay returns all records for a specific day for the authenticated user.
func (s *Service) ListByDay(ctx context.Context, userID uint64, date time.Time) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListByDay)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListByDay),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("date", date.Format("2006-01-02")),
	)

	span.AddEvent("CheckCache")
	cachedRecords, err := s.RecordCache.GetRecordsByDay(ctx, userID, date)
	if err == nil && cachedRecords != nil {
		span.AddEvent("CacheHit")
		span.SetStatus(codes.Ok, StatusListedAll)
		s.Logger.InfowCtx(ctx, "records retrieved from cache",
			commonkeys.UserID, userID,
			"date", date.Format("2006-01-02"),
			"count", len(cachedRecords),
		)
		return cachedRecords, nil
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListByDay(ctx, userID, date)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.UserID, userID,
			"date", date.Format("2006-01-02"),
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent("SaveToCache")
	if err := s.RecordCache.SaveRecordsByDay(ctx, userID, date, records, 0); err != nil {
		s.Logger.WarnwCtx(ctx, "failed to save records to cache",
			commonkeys.UserID, userID,
			"date", date.Format("2006-01-02"),
			commonkeys.Error, err,
		)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed by day successfully",
		commonkeys.UserID, userID,
		"date", date.Format("2006-01-02"),
		"count", len(records),
	)

	return records, nil
}
