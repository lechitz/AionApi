package usecase

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/audit/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// WriteEvent validates and persists one immutable audit event.
func (s *Service) WriteEvent(ctx context.Context, event domain.AuditActionEvent) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanWriteEvent)
	defer span.End()

	event = normalizeEvent(event)
	span.SetAttributes(
		attribute.String(LogKeyUserID, strconv.FormatUint(event.UserID, 10)),
		attribute.String(LogKeyTraceID, event.TraceID),
		attribute.String(LogKeyDraftID, event.DraftID),
		attribute.String("status", event.Status),
	)

	s.logger.InfowCtx(ctx, LogWritingAuditEvent,
		LogKeyUserID, event.UserID,
		LogKeyTraceID, event.TraceID,
		LogKeyDraftID, event.DraftID,
	)

	if err := validateEvent(event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	if err := s.repository.Save(ctx, event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedWriteAuditEvent)
		s.logger.ErrorwCtx(ctx, LogFailedWriteAuditEvent,
			LogKeyError, err.Error(),
			LogKeyUserID, event.UserID,
			LogKeyDraftID, event.DraftID,
		)
		return err
	}

	span.SetStatus(codes.Ok, StatusAuditEventWritten)
	s.logger.InfowCtx(ctx, LogAuditEventWritten,
		LogKeyUserID, event.UserID,
		LogKeyTraceID, event.TraceID,
		LogKeyDraftID, event.DraftID,
	)
	return nil
}

func normalizeEvent(event domain.AuditActionEvent) domain.AuditActionEvent {
	if strings.TrimSpace(event.EventID) == "" {
		event.EventID = uuid.NewString()
	}
	if event.TimestampUTC.IsZero() {
		event.TimestampUTC = time.Now().UTC()
	}
	if strings.TrimSpace(event.Source) == "" {
		event.Source = "aion-api"
	}
	return event
}

func validateEvent(event domain.AuditActionEvent) error {
	if event.UserID == 0 {
		return errors.New("audit event requires user_id")
	}
	if strings.TrimSpace(event.EventID) == "" {
		return errors.New("audit event requires event_id")
	}
	if event.TimestampUTC.IsZero() {
		return errors.New("audit event requires timestamp_utc")
	}
	if strings.TrimSpace(event.UIActionType) == "" {
		return errors.New("audit event requires ui_action_type")
	}
	if strings.TrimSpace(event.DraftID) == "" {
		return errors.New("audit event requires draft_id")
	}
	if strings.TrimSpace(event.Status) == "" {
		return errors.New("audit event requires status")
	}
	return nil
}
