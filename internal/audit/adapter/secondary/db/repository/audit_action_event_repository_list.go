package repository

import (
	"context"

	"github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// List retrieves audit events by optional filters for internal diagnostics.
func (r *AuditActionEventRepository) List(ctx context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error) {
	tr := otel.Tracer(AuditTracerName)
	ctx, span := tr.Start(ctx, SpanAuditListRepo, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpAuditList),
		attribute.String("trace_id", filter.TraceID),
		attribute.String("draft_id", filter.DraftID),
		attribute.Int("statuses_count", len(filter.Statuses)),
		attribute.Int("limit", filter.Limit),
		attribute.Int("offset", filter.Offset),
	))
	defer span.End()

	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}

	query := r.db.WithContext(ctx).Model(&model.AuditActionEventDB{})

	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.TraceID != "" {
		query = query.Where("trace_id = ?", filter.TraceID)
	}
	if filter.DraftID != "" {
		query = query.Where("draft_id = ?", filter.DraftID)
	}
	if len(filter.Statuses) > 0 {
		query = query.Where("status IN ?", filter.Statuses)
	}
	if filter.FromUTC != nil {
		query = query.Where("timestamp_utc >= ?", *filter.FromUTC)
	}
	if filter.ToUTC != nil {
		query = query.Where("timestamp_utc <= ?", *filter.ToUTC)
	}

	var rows []model.AuditActionEventDB
	err := query.
		Order("timestamp_utc DESC").
		Limit(limit).
		Offset(filter.Offset).
		Find(&rows).
		Error()
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpAuditList)
		r.logger.ErrorwCtx(ctx, ErrListAuditActionEventMsg,
			commonkeys.Error, err.Error(),
			"trace_id", filter.TraceID,
			"draft_id", filter.DraftID,
			"limit", limit,
		)
		return nil, err
	}

	events := mapper.AuditActionEventsFromDB(rows)
	span.SetAttributes(attribute.Int("results_count", len(events)))
	span.SetStatus(codes.Ok, StatusAuditListed)

	r.logger.InfowCtx(ctx, StatusAuditListed,
		"results_count", len(events),
		"trace_id", filter.TraceID,
		"draft_id", filter.DraftID,
		"limit", limit,
	)

	return events, nil
}
