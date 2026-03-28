package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
	"github.com/lechitz/aion-api/internal/shared/constants/ctxkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Enqueue validates and persists one canonical event in the durable outbox.
func (s *Service) Enqueue(ctx context.Context, event domain.Event) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanEnqueue)
	defer span.End()

	span.AddEvent(EventNormalizeInput)
	event = s.normalizeEvent(ctx, event)

	span.SetAttributes(
		attribute.String(LogKeyEventID, event.EventID),
		attribute.String(LogKeyAggregateType, event.AggregateType),
		attribute.String(LogKeyAggregateID, event.AggregateID),
		attribute.String(LogKeyEventType, event.EventType),
		attribute.String(LogKeyEventVersion, event.EventVersion),
	)

	s.logger.InfowCtx(ctx, LogEnqueueingEvent,
		LogKeyEventID, event.EventID,
		LogKeyAggregateType, event.AggregateType,
		LogKeyAggregateID, event.AggregateID,
		LogKeyEventType, event.EventType,
		LogKeyEventVersion, event.EventVersion,
	)

	span.AddEvent(EventValidateInput)
	if err := validateEvent(event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}

	span.AddEvent(EventRepositorySave)
	if err := s.repository.Save(ctx, event); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, LogFailedToEnqueueEvent)
		s.logger.ErrorwCtx(ctx, LogFailedToEnqueueEvent,
			"error", err.Error(),
			LogKeyEventID, event.EventID,
			LogKeyEventType, event.EventType,
			LogKeyAggregateType, event.AggregateType,
			LogKeyAggregateID, event.AggregateID,
		)
		return err
	}

	span.AddEvent(EventSuccess)
	span.SetStatus(codes.Ok, StatusEventQueued)
	s.logger.InfowCtx(ctx, LogOutboxEventQueued,
		LogKeyEventID, event.EventID,
		LogKeyEventType, event.EventType,
		LogKeyAggregateType, event.AggregateType,
		LogKeyAggregateID, event.AggregateID,
	)
	return nil
}

func (s *Service) normalizeEvent(ctx context.Context, event domain.Event) domain.Event {
	if strings.TrimSpace(event.EventID) == "" {
		event.EventID = uuid.NewString()
	}
	if strings.TrimSpace(event.EventVersion) == "" {
		event.EventVersion = EventVersionV1
	}
	if strings.TrimSpace(event.Source) == "" {
		event.Source = DefaultSource
	}
	if strings.TrimSpace(event.Status) == "" {
		event.Status = DefaultEventStatus
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = s.now()
	}
	if event.AvailableAtUTC.IsZero() {
		event.AvailableAtUTC = event.CreatedAt
	}
	if strings.TrimSpace(event.TraceID) == "" {
		event.TraceID = traceIDFromContext(ctx)
	}
	if strings.TrimSpace(event.RequestID) == "" {
		if requestID, ok := ctx.Value(ctxkeys.RequestID).(string); ok {
			event.RequestID = requestID
		}
	}
	return event
}

func validateEvent(event domain.Event) error {
	switch {
	case strings.TrimSpace(event.EventID) == "":
		return ErrEventIDRequired
	case strings.TrimSpace(event.AggregateType) == "":
		return ErrAggregateTypeRequired
	case strings.TrimSpace(event.AggregateID) == "":
		return ErrAggregateIDRequired
	case strings.TrimSpace(event.EventType) == "":
		return ErrEventTypeRequired
	case strings.TrimSpace(event.EventVersion) == "":
		return ErrEventVersionRequired
	case strings.TrimSpace(event.Source) == "":
		return ErrSourceRequired
	case len(event.PayloadJSON) == 0:
		return ErrPayloadRequired
	default:
		return nil
	}
}

func traceIDFromContext(ctx context.Context) string {
	spanContext := trace.SpanFromContext(ctx).SpanContext()
	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}

	if traceID, ok := ctx.Value(ctxkeys.TraceID).(string); ok {
		return traceID
	}
	return ""
}
