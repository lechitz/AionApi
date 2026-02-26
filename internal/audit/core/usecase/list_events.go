package usecase

import (
	"context"

	"github.com/lechitz/AionApi/internal/audit/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// ListEvents returns audit events according to filters.
func (s *Service) ListEvents(ctx context.Context, filter domain.AuditActionEventFilter) ([]domain.AuditActionEvent, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanListEvents)
	defer span.End()

	if filter.Limit <= 0 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	span.SetAttributes(
		attribute.String(LogKeyTraceID, filter.TraceID),
		attribute.String(LogKeyDraftID, filter.DraftID),
		attribute.Int("limit", filter.Limit),
		attribute.Int("offset", filter.Offset),
	)

	s.logger.InfowCtx(ctx, LogListingAuditEvents,
		LogKeyTraceID, filter.TraceID,
		LogKeyDraftID, filter.DraftID,
		"limit", filter.Limit,
	)

	events, err := s.repository.List(ctx, filter)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedListAuditEvents)
		s.logger.ErrorwCtx(ctx, LogFailedListAuditEvents,
			LogKeyError, err.Error(),
			LogKeyTraceID, filter.TraceID,
			LogKeyDraftID, filter.DraftID,
		)
		return nil, err
	}

	span.SetAttributes(attribute.Int("results_count", len(events)))
	span.SetStatus(codes.Ok, StatusAuditEventsListed)
	return events, nil
}
