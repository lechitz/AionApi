package repository

import (
	"context"
	"fmt"

	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/eventoutbox/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// ListByStatus returns outbox events filtered by status for operational diagnostics.
func (r *EventRepository) ListByStatus(ctx context.Context, status string, limit int) ([]domain.Event, error) {
	var rows []model.EventDB

	query := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("available_at_utc ASC, id ASC").
		Limit(limit).
		Find(&rows)
	if err := query.Error(); err != nil {
		r.logger.ErrorwCtx(ctx, "error listing outbox events by status",
			commonkeys.Error, err.Error(),
			"status", status,
			"limit", limit,
		)
		return nil, fmt.Errorf("list outbox events by status: %w", err)
	}

	return mapper.EventsFromDB(rows), nil
}

// ListFailed returns pending outbox events that already have a last_error value.
func (r *EventRepository) ListFailed(ctx context.Context, limit int) ([]domain.Event, error) {
	var rows []model.EventDB

	query := r.db.WithContext(ctx).
		Where("status = ? AND last_error IS NOT NULL AND last_error <> ''", "pending").
		Order("available_at_utc ASC, id ASC").
		Limit(limit).
		Find(&rows)
	if err := query.Error(); err != nil {
		r.logger.ErrorwCtx(ctx, "error listing failed outbox events",
			commonkeys.Error, err.Error(),
			"limit", limit,
		)
		return nil, fmt.Errorf("list failed outbox events: %w", err)
	}

	return mapper.EventsFromDB(rows), nil
}
