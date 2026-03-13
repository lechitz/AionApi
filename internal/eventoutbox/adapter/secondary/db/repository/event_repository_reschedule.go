package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Reschedule increments attempts and postpones the next publish attempt.
func (r *EventRepository) Reschedule(ctx context.Context, eventID string, nextAvailableAt time.Time, lastError string) error {
	tr := otel.Tracer(OutboxTracerName)
	ctx, span := tr.Start(ctx, SpanOutboxRescheduleRepo, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpOutboxReschedule),
		attribute.String("event_id", eventID),
	))
	defer span.End()

	res := r.db.WithContext(ctx).Exec(
		`UPDATE aion_api.event_outbox
		 SET attempt_count = attempt_count + 1,
		     available_at_utc = ?,
		     last_error = ?,
		     updated_at = ?
		 WHERE event_id = ?`,
		nextAvailableAt,
		lastError,
		time.Now().UTC(),
		eventID,
	)

	if err := res.Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxReschedule)
		r.logger.ErrorwCtx(ctx, ErrRescheduleOutboxEventMsg,
			commonkeys.Error, err.Error(),
			"event_id", eventID,
		)
		return fmt.Errorf("reschedule outbox event: %w", err)
	}

	span.SetStatus(codes.Ok, StatusOutboxRescheduled)
	return nil
}
