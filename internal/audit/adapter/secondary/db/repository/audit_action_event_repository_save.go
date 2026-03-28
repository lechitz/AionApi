package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/audit/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Save inserts one immutable audit action event row.
func (r *AuditActionEventRepository) Save(ctx context.Context, event domain.AuditActionEvent) error {
	tr := otel.Tracer(AuditTracerName)
	ctx, span := tr.Start(ctx, SpanAuditSaveRepo, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(event.UserID, 10)),
		attribute.String(commonkeys.Operation, OpAuditSave),
		attribute.String("event_id", event.EventID),
		attribute.String("ui_action_type", event.UIActionType),
		attribute.String("draft_id", event.DraftID),
		attribute.String(commonkeys.Status, event.Status),
	))
	defer span.End()

	row := mapper.AuditActionEventToDB(event)

	if err := r.db.WithContext(ctx).Create(&row).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpAuditSave)
		r.logger.ErrorwCtx(ctx, ErrSaveAuditActionEventMsg,
			commonkeys.Error, err.Error(),
			commonkeys.UserID, strconv.FormatUint(event.UserID, 10),
			"event_id", event.EventID,
			"draft_id", event.DraftID,
		)
		return fmt.Errorf("save audit action event: %w", err)
	}

	span.SetStatus(codes.Ok, StatusAuditSaved)
	r.logger.InfowCtx(ctx, StatusAuditSaved,
		commonkeys.UserID, strconv.FormatUint(event.UserID, 10),
		"event_id", event.EventID,
		"draft_id", event.DraftID,
	)
	return nil
}
