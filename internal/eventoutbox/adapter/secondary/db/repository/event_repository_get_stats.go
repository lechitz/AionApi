package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type outboxStatsRow struct {
	PendingCount       int64      `gorm:"column:pending_count"`
	PublishedCount     int64      `gorm:"column:published_count"`
	FailedCount        int64      `gorm:"column:failed_count"`
	OldestPendingAtUTC *time.Time `gorm:"column:oldest_pending_at_utc"`
}

const getOutboxStatsSQL = `
SELECT
	COUNT(*) FILTER (WHERE status = 'pending') AS pending_count,
	COUNT(*) FILTER (WHERE status = 'published') AS published_count,
	COUNT(*) FILTER (WHERE status = 'pending' AND last_error IS NOT NULL AND last_error <> '') AS failed_count,
	MIN(available_at_utc) FILTER (WHERE status = 'pending') AS oldest_pending_at_utc
FROM aion_api.event_outbox;
`

func (r *EventRepository) GetStats(ctx context.Context) (domain.Stats, error) {
	tr := otel.Tracer(OutboxTracerName)
	ctx, span := tr.Start(ctx, "eventoutbox.repository.get_stats", trace.WithAttributes(
		attribute.String(commonkeys.Operation, "event_outbox_get_stats"),
	))
	defer span.End()

	var row outboxStatsRow
	query := r.db.WithContext(ctx).Raw(getOutboxStatsSQL).Scan(&row)
	if err := query.Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "event_outbox_get_stats")
		r.logger.ErrorwCtx(ctx, "error reading outbox stats", commonkeys.Error, err.Error())
		return domain.Stats{}, fmt.Errorf("get outbox stats: %w", err)
	}

	stats := domain.Stats{
		PendingCount:       row.PendingCount,
		PublishedCount:     row.PublishedCount,
		FailedCount:        row.FailedCount,
		OldestPendingAtUTC: row.OldestPendingAtUTC,
	}

	span.SetAttributes(
		attribute.Int64("pending_count", stats.PendingCount),
		attribute.Int64("published_count", stats.PublishedCount),
		attribute.Int64("failed_count", stats.FailedCount),
	)
	span.SetStatus(codes.Ok, "event outbox stats loaded")
	return stats, nil
}
