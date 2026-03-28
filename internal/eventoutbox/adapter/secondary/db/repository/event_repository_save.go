package repository

import (
	"context"
	"fmt"

	"github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Save inserts one durable outbox event row.
func (r *EventRepository) Save(ctx context.Context, event domain.Event) error {
	tr := otel.Tracer(OutboxTracerName)
	ctx, span := tr.Start(ctx, SpanOutboxSaveRepo, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpOutboxSave),
		attribute.String(commonkeys.Entity, event.AggregateType),
		attribute.String("aggregate_id", event.AggregateID),
		attribute.String("event_type", event.EventType),
		attribute.String("event_id", event.EventID),
	))
	defer span.End()

	row := mapper.EventToDB(event)

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxSave)
		r.logger.ErrorwCtx(ctx, ErrSaveOutboxEventMsg,
			commonkeys.Error, err.Error(),
			commonkeys.Entity, event.AggregateType,
			"aggregate_id", event.AggregateID,
			"event_type", event.EventType,
			"event_id", event.EventID,
		)
		return fmt.Errorf("save outbox event: %w", err)
	}

	span.SetStatus(codes.Ok, StatusOutboxSaved)
	r.logger.InfowCtx(ctx, StatusOutboxSaved,
		commonkeys.Entity, event.AggregateType,
		"aggregate_id", event.AggregateID,
		"event_type", event.EventType,
		"event_id", event.EventID,
	)

	return nil
}
