package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// MarkPublished marks one outbox event as published.
func (r *EventRepository) MarkPublished(ctx context.Context, eventID string, publishedAt time.Time) error {
	tr := otel.Tracer(OutboxTracerName)
	ctx, span := tr.Start(ctx, SpanOutboxMarkPublishedRepo, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpOutboxMarkPublished),
		attribute.String("event_id", eventID),
	))
	defer span.End()

	res := r.db.WithContext(ctx).
		Model(&model.EventDB{}).
		Where("event_id = ?", eventID).
		Updates(map[string]any{
			"status":           "published",
			"published_at_utc": publishedAt,
			"last_error":       "",
			"updated_at":       publishedAt,
		})

	if err := res.Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxMarkPublished)
		r.logger.ErrorwCtx(ctx, ErrMarkOutboxPublishedMsg,
			commonkeys.Error, err.Error(),
			"event_id", eventID,
		)
		return fmt.Errorf("mark outbox published: %w", err)
	}

	span.SetStatus(codes.Ok, StatusOutboxPublished)
	return nil
}
