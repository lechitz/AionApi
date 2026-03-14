package usecase

import (
	"context"
	"fmt"
	"strconv"

	eventoutboxinput "github.com/lechitz/AionApi/internal/eventoutbox/core/ports/input"
	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/lechitz/AionApi/internal/record/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Update applies partial changes to an existing record owned by the user.
func (s *Service) Update(ctx context.Context, recordID uint64, userID uint64, cmd input.UpdateRecordCommand) (domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanUpdate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdate),
		attribute.String(commonkeys.RecordID, strconv.FormatUint(recordID, 10)),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	span.AddEvent(EventValidateInput)
	if recordID == 0 || userID == 0 {
		span.RecordError(ErrInvalidRecordIDOrUserID)
		span.SetStatus(codes.Error, ErrToValidateRecord)
		s.Logger.ErrorwCtx(ctx, ErrToValidateRecord, commonkeys.Error, ErrInvalidRecordIDOrUserID.Error())
		return domain.Record{}, ErrInvalidRecordIDOrUserID
	}

	if cmd.TagID != nil {
		tag, err := s.TagRepository.GetByID(ctx, *cmd.TagID, userID)
		if err != nil || tag.ID == 0 {
			span.RecordError(err)
			span.SetStatus(codes.Error, FailedToUpdateRecord)
			s.Logger.ErrorwCtx(ctx, FailedToUpdateRecord, commonkeys.Error, TagNotFound)
			return domain.Record{}, fmt.Errorf("%w: %s", ErrUpdateRecord, TagNotFound)
		}
	}

	var updated domain.Record
	if err := s.runWithinRecordOutboxTransaction(ctx, func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		span.AddEvent(EventRepositoryGet)
		existing, getErr := recordRepo.GetByID(ctx, recordID, userID)
		if getErr != nil {
			return fmt.Errorf("%w: %w", ErrGetRecord, getErr)
		}

		finalTagID := existing.TagID
		if cmd.TagID != nil {
			finalTagID = *cmd.TagID
		}

		existing = applyRecordPatch(existing, cmd, finalTagID)

		span.AddEvent(EventRepositoryUpdate)
		var updateErr error
		updated, updateErr = recordRepo.Update(ctx, existing)
		if updateErr != nil {
			return updateErr
		}

		if outboxService != nil {
			s.enqueueRecordOutboxEventWithService(ctx, outboxService, RecordEventTypeUpdatedV1, updated)
		}

		return nil
	}); err != nil {
		span.RecordError(err)
		if isGetRecordError(err) {
			span.SetStatus(codes.Error, FailedToGetRecord)
			s.Logger.ErrorwCtx(ctx, FailedToGetRecord,
				commonkeys.RecordID, recordID,
				commonkeys.UserID, userID,
				commonkeys.Error, err,
			)
			return domain.Record{}, err
		}
		span.SetStatus(codes.Error, FailedToUpdateRecord)
		s.Logger.ErrorwCtx(ctx, FailedToUpdateRecord,
			commonkeys.RecordID, recordID,
			commonkeys.Error, err,
		)
		return domain.Record{}, fmt.Errorf("%w: %w", ErrUpdateRecord, err)
	}

	// Invalidate all related caches
	s.invalidateRecordCaches(ctx, span, updated)

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusUpdated)
	s.Logger.InfowCtx(ctx, LogRecordUpdatedSuccessfully,
		commonkeys.RecordID, updated.ID,
		commonkeys.UserID, updated.UserID,
	)

	return updated, nil
}

// applyRecordPatch mutates a copy of the record with fields from cmd and the resolved tag ID.
func applyRecordPatch(r domain.Record, cmd input.UpdateRecordCommand, tagID uint64) domain.Record {
	if cmd.Description != nil {
		r.Description = cmd.Description
	}
	r.TagID = tagID
	if cmd.EventTime != nil {
		r.EventTime = *cmd.EventTime
	}
	if cmd.RecordedAt != nil {
		r.RecordedAt = cmd.RecordedAt
	}
	if cmd.DurationSecs != nil {
		r.DurationSecs = cmd.DurationSecs
	}
	if cmd.Value != nil {
		r.Value = cmd.Value
	}
	if cmd.Source != nil {
		r.Source = cmd.Source
	}
	if cmd.Timezone != nil {
		r.Timezone = cmd.Timezone
	}
	if cmd.Status != nil {
		r.Status = cmd.Status
	}
	return r
}

// invalidateRecordCaches invalidates all caches related to the updated record.
// This is a best-effort operation - errors are logged but don't fail the operation.
func (s *Service) invalidateRecordCaches(ctx context.Context, span trace.Span, record domain.Record) {
	span.AddEvent(EventInvalidateCache)

	// Invalidate the specific record cache
	if err := s.RecordCache.DeleteRecord(ctx, record.ID, record.UserID); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedInvalidateRecordCache,
			commonkeys.RecordID, record.ID,
			commonkeys.UserID, record.UserID,
			commonkeys.Error, err,
		)
	}

	// Invalidate day cache for the event date
	eventDate := CacheDayStart(record.EventTime)
	if err := s.RecordCache.DeleteRecordsByDay(ctx, record.UserID, eventDate); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedInvalidateDayCache,
			commonkeys.UserID, record.UserID,
			commonkeys.Date, eventDate.Format(DateFormatISO8601Date),
			commonkeys.Error, err,
		)
	}

	// Get tag to find category for cache invalidation
	if record.TagID != 0 {
		tag, err := s.TagRepository.GetByID(ctx, record.TagID, record.UserID)
		if err == nil && tag.ID != 0 {
			// Invalidate category cache
			if err := s.RecordCache.DeleteRecordsByCategory(ctx, tag.CategoryID, record.UserID); err != nil {
				s.Logger.WarnwCtx(ctx, LogFailedInvalidateCategoryCache,
					commonkeys.CategoryID, tag.CategoryID,
					commonkeys.UserID, record.UserID,
					commonkeys.Error, err,
				)
			}
		}

		// Invalidate tag cache
		if err := s.RecordCache.DeleteRecordsByTag(ctx, record.TagID, record.UserID); err != nil {
			s.Logger.WarnwCtx(ctx, LogFailedInvalidateTagCache,
				commonkeys.TagID, record.TagID,
				commonkeys.UserID, record.UserID,
				commonkeys.Error, err,
			)
		}
	}
}
