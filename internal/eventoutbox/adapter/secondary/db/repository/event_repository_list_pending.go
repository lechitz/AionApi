package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ListPending returns pending outbox rows available for publication.
func (r *EventRepository) ListPending(ctx context.Context, limit int) ([]domain.Event, error) {
	tr := otel.Tracer(OutboxTracerName)
	ctx, span := tr.Start(ctx, SpanOutboxListPendingRepo, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpOutboxListPending),
		attribute.Int("limit", limit),
	))
	defer span.End()

	var rows []model.EventDB
	query := r.db.WithContext(ctx).
		Where("status = ? AND available_at_utc <= ?", "pending", time.Now().UTC()).
		Order("available_at_utc ASC, id ASC").
		Limit(limit).
		Find(&rows)

	if err := query.Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxListPending)
		r.logger.ErrorwCtx(ctx, ErrListOutboxEventsMsg,
			commonkeys.Error, err.Error(),
			"limit", limit,
		)
		return nil, fmt.Errorf("list pending outbox events: %w", err)
	}

	span.SetAttributes(attribute.Int("rows", len(rows)))
	span.SetStatus(codes.Ok, StatusOutboxListed)
	return mapper.EventsFromDB(rows), nil
}
