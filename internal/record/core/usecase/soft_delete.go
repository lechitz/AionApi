package usecase

import (
	"context"
	"fmt"
	"strconv"

	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Delete performs a soft delete for the user's record.
func (s *Service) Delete(ctx context.Context, id uint64, userID uint64) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanSoftDelete)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanSoftDelete),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(id, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	var existing domain.Record
	if err := s.runWithinRecordOutboxTransaction(ctx, func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		span.AddEvent(EventRepositoryGet)
		var getErr error
		existing, getErr = recordRepo.GetByID(ctx, id, userID)
		if getErr != nil {
			return fmt.Errorf("%w: %w", ErrGetRecord, getErr)
		}

		span.AddEvent(EventRepositoryDelete)
		if deleteErr := recordRepo.Delete(ctx, id, userID); deleteErr != nil {
			return deleteErr
		}

		if outboxService != nil {
			s.enqueueRecordOutboxEventWithService(ctx, outboxService, RecordEventTypeDeletedV1, existing)
		}
		return nil
	}); err != nil {
		span.RecordError(err)
		if isGetRecordError(err) {
			span.SetStatus(codes.Error, FailedToGetRecord)
			s.Logger.ErrorwCtx(ctx, FailedToGetRecord,
				commonkeys.RecordID, id,
				commonkeys.UserID, userID,
				commonkeys.Error, err,
			)
			return err
		}
		span.SetStatus(codes.Error, FailedToDeleteRecord)
		s.Logger.ErrorwCtx(ctx, FailedToDeleteRecord,
			commonkeys.RecordID, id,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
		return fmt.Errorf("%w: %w", ErrDeleteRecord, err)
	}

	span.AddEvent(EventInvalidateCache)
	if err := s.RecordCache.DeleteRecord(ctx, id, userID); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedInvalidateRecordCache,
			commonkeys.RecordID, id,
			commonkeys.UserID, userID,
			commonkeys.Error, err,
		)
	}

	// Invalidate day cache
	eventDate := CacheDayStart(existing.EventTime)
	if err := s.RecordCache.DeleteRecordsByDay(ctx, userID, eventDate); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedInvalidateDayCache,
			commonkeys.UserID, userID,
			commonkeys.Date, eventDate.Format(DateFormatISO8601Date),
			commonkeys.Error, err,
		)
	}

	// Get tag to find category for cache invalidation
	if existing.TagID != 0 {
		tagObj, err := s.TagRepository.GetByID(ctx, existing.TagID, userID)
		if err == nil && tagObj.ID != 0 {
			if err := s.RecordCache.DeleteRecordsByCategory(ctx, tagObj.CategoryID, userID); err != nil {
				s.Logger.WarnwCtx(ctx, LogFailedInvalidateCategoryCache,
					commonkeys.CategoryID, tagObj.CategoryID,
					commonkeys.UserID, userID,
					commonkeys.Error, err,
				)
			}
		}

		if err := s.RecordCache.DeleteRecordsByTag(ctx, existing.TagID, userID); err != nil {
			s.Logger.WarnwCtx(ctx, LogFailedInvalidateTagCache,
				commonkeys.TagID, existing.TagID,
				commonkeys.UserID, userID,
				commonkeys.Error, err,
			)
		}
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusDeleted)
	s.Logger.InfowCtx(ctx, LogRecordSoftDeletedSuccess,
		commonkeys.RecordID, id,
		commonkeys.UserID, userID,
	)

	return nil
}
