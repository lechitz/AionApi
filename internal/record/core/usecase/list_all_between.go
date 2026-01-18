package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListAllBetween returns records with event_time within the specified date range (inclusive).
func (s *Service) ListAllBetween(ctx context.Context, userID uint64, startDate time.Time, endDate time.Time, limit int) ([]domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListAllBetween)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanListAllBetween),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String("start_date", startDate.Format("2006-01-02")),
		attribute.String("end_date", endDate.Format("2006-01-02")),
		attribute.Int("limit", limit),
	)

	if limit <= 0 || limit > 100 {
		limit = 50 // default limit
	}

	// Validate date range
	span.AddEvent(EventValidateInput)
	if startDate.After(endDate) {
		err := errors.New(StartDateMustBeBeforeEndDate)
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateRecord)
		s.Logger.ErrorwCtx(ctx, ErrToValidateRecord, commonkeys.Error, err.Error())
		return nil, err
	}

	span.AddEvent(EventRepositoryList)
	records, err := s.RecordRepository.ListAllBetween(ctx, userID, startDate, endDate, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToListRecords)
		s.Logger.ErrorwCtx(ctx, FailedToListRecords,
			commonkeys.UserID, userID,
			"start_date", startDate.Format("2006-01-02"),
			"end_date", endDate.Format("2006-01-02"),
			commonkeys.Error, err,
		)
		return nil, fmt.Errorf("%s: %w", FailedToListRecords, err)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusListedAll)
	s.Logger.InfowCtx(ctx, "records listed between dates successfully",
		commonkeys.UserID, userID,
		"start_date", startDate.Format("2006-01-02"),
		"end_date", endDate.Format("2006-01-02"),
		"count", len(records),
	)

	return records, nil
}
