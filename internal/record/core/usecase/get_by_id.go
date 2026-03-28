package usecase

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetByID returns a record by id for the user.
func (s *Service) GetByID(ctx context.Context, recordID uint64, userID uint64) (domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanGetByID)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanGetByID),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(recordID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	span.AddEvent("CheckCache")
	cachedRecord, err := s.RecordCache.GetRecord(ctx, recordID, userID)
	if err == nil && cachedRecord.ID != 0 {
		span.AddEvent("CacheHit")
		span.SetStatus(codes.Ok, StatusRetrieved)
		s.Logger.InfowCtx(ctx, "record retrieved from cache",
			commonkeys.RecordID, recordID,
			commonkeys.UserID, userID,
		)
		return cachedRecord, nil
	}

	span.AddEvent(EventRepositoryGet)
	rec, err := s.RecordRepository.GetByID(ctx, recordID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToGetRecord)
		s.Logger.ErrorwCtx(ctx, FailedToGetRecord,
			commonkeys.RecordID, recordID,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return domain.Record{}, fmt.Errorf("%w: %w", ErrGetRecord, err)
	}

	span.AddEvent("SaveToCache")
	if err := s.RecordCache.SaveRecord(ctx, rec, 0); err != nil {
		s.Logger.WarnwCtx(ctx, "failed to save record to cache",
			commonkeys.RecordID, rec.ID,
			commonkeys.UserID, rec.UserID,
			commonkeys.Error, err,
		)
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusRetrieved)
	s.Logger.InfowCtx(ctx, "record retrieved from database",
		commonkeys.RecordID, rec.ID,
		commonkeys.UserID, rec.UserID,
	)

	return rec, nil
}
