package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

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

// Create creates a new record after validating inputs.
func (s *Service) Create(ctx context.Context, cmd input.CreateRecordCommand) (domain.Record, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanCreate),
	)

	span.AddEvent(EventValidateInput)

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateRecord)
		s.Logger.ErrorwCtx(ctx, ErrToValidateRecord, commonkeys.Error, err.Error())
		return domain.Record{}, err
	}

	eventTime := resolveEventTime(cmd)

	if err := validateRecordedAt(cmd.RecordedAt); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrToValidateRecord)
		s.Logger.ErrorwCtx(ctx, ErrToValidateRecord, commonkeys.Error, err.Error())
		return domain.Record{}, err
	}

	finalTagID, err := s.resolveTagID(ctx, cmd.TagID, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToCreateRecord)
		s.Logger.ErrorwCtx(ctx, FailedToCreateRecord, commonkeys.Error, err.Error())
		return domain.Record{}, fmt.Errorf("%w: %w", ErrCreateRecord, err)
	}

	// Apply default values for optional fields
	recordedAt := resolveRecordedAt(cmd.RecordedAt)
	status := resolveStatus(cmd.Status)
	timezone := resolveTimezone(cmd.Timezone)

	rec := domain.Record{
		UserID:       userID,
		Description:  cmd.Description,
		TagID:        finalTagID,
		EventTime:    eventTime,
		RecordedAt:   recordedAt,
		DurationSecs: cmd.DurationSecs,
		Value:        cmd.Value,
		Source:       cmd.Source,
		Timezone:     timezone,
		Status:       status,
	}

	var created domain.Record
	if err := s.runWithinRecordOutboxTransaction(ctx, func(recordRepo output.RecordRepository, outboxService eventoutboxinput.Service) error {
		span.AddEvent(EventRepositoryCreate)
		var createErr error
		created, createErr = recordRepo.Create(ctx, rec)
		if createErr != nil {
			return createErr
		}

		if outboxService != nil {
			s.enqueueRecordOutboxEventWithService(ctx, outboxService, RecordEventTypeCreatedV1, created)
		}
		return nil
	}); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, FailedToCreateRecord)
		s.Logger.ErrorwCtx(ctx, FailedToCreateRecord, commonkeys.Error, err)
		return domain.Record{}, fmt.Errorf("%w: %w", ErrCreateRecord, err)
	}

	s.saveToCacheAndInvalidate(ctx, span, created)

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusCreated)
	s.Logger.InfowCtx(ctx, LogRecordCreatedSuccessfully,
		commonkeys.RecordID, created.ID,
		commonkeys.UserID, created.UserID,
	)

	return created, nil
}

// resolveEventTime determines the event time from the command.
func resolveEventTime(cmd input.CreateRecordCommand) time.Time {
	if !cmd.EventTime.IsZero() {
		return cmd.EventTime
	}
	if cmd.RecordedAt != nil {
		return *cmd.RecordedAt
	}
	return time.Now().UTC()
}

// validateRecordedAt ensures recordedAt is not in the future.
func validateRecordedAt(recordedAt *time.Time) error {
	if recordedAt != nil && recordedAt.After(time.Now().UTC()) {
		return ErrRecordedAtFuture
	}
	return nil
}

// resolveRecordedAt returns the provided recordedAt or defaults to now.
func resolveRecordedAt(recordedAt *time.Time) *time.Time {
	if recordedAt != nil {
		return recordedAt
	}
	now := time.Now().UTC()
	return &now
}

// resolveStatus returns the provided status or defaults to "published".
func resolveStatus(status *string) *string {
	if status != nil && *status != "" {
		return status
	}
	defaultStatus := DefaultRecordStatus
	return &defaultStatus
}

// resolveTimezone returns the provided timezone or defaults to "America/Sao_Paulo".
func resolveTimezone(timezone *string) *string {
	if timezone != nil && *timezone != "" {
		return timezone
	}
	defaultTZ := DefaultTimezone
	return &defaultTZ
}

// resolveTagID validates that the tag exists and belongs to the user.
func (s *Service) resolveTagID(ctx context.Context, tagID uint64, userID uint64) (uint64, error) {
	if tagID == 0 {
		return 0, ErrTagIDIsRequired
	}

	// Validate tag exists and belongs to user
	tagObj, err := s.TagRepository.GetByID(ctx, tagID, userID)
	if err != nil {
		return 0, fmt.Errorf(ErrLookupTagFormat, err)
	}

	if tagObj.ID == 0 {
		return 0, errors.New(TagNotFound)
	}

	return tagID, nil
}

// saveToCacheAndInvalidate saves the record to cache and invalidates related list caches.
// This is a best-effort operation - errors are logged but don't fail the operation.
func (s *Service) saveToCacheAndInvalidate(ctx context.Context, span trace.Span, record domain.Record) {
	span.AddEvent(EventSaveToCache)
	if err := s.RecordCache.SaveRecord(ctx, record, 0); err != nil {
		s.Logger.WarnwCtx(ctx, LogFailedSaveRecordToCacheAfterCreation,
			commonkeys.RecordID, record.ID,
			commonkeys.UserID, record.UserID,
			commonkeys.Error, err,
		)
	}

	span.AddEvent(EventInvalidateCache)
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
		tagObj, err := s.TagRepository.GetByID(ctx, record.TagID, record.UserID)
		if err == nil && tagObj.ID != 0 {
			// Invalidate category cache (records are accessed via tag → category)
			if err := s.RecordCache.DeleteRecordsByCategory(ctx, tagObj.CategoryID, record.UserID); err != nil {
				s.Logger.WarnwCtx(ctx, LogFailedInvalidateCategoryCache,
					commonkeys.CategoryID, tagObj.CategoryID,
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
