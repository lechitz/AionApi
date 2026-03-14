package usecase

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// PublishPending publishes one batch of pending outbox rows.
func (s *PublisherService) PublishPending(ctx context.Context, limit int) error {
	if limit <= 0 {
		limit = 1
	}

	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanPublishPending)
	defer span.End()

	span.SetAttributes(attribute.Int("batch_limit", limit))
	s.logger.InfowCtx(ctx, LogPublishPendingEvents, "limit", limit)

	events, err := s.repository.ListPending(ctx, limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, EventRepositoryList)
		return err
	}

	var publishErr error
	for _, event := range events {
		if err := s.publisher.Publish(ctx, event); err != nil {
			s.logger.ErrorwCtx(ctx, LogOutboxEventPublishFailed,
				commonkeys.Error, err.Error(),
				LogKeyEventID, event.EventID,
				LogKeyEventType, event.EventType,
				LogKeyAggregateType, event.AggregateType,
				LogKeyAggregateID, event.AggregateID,
			)

			if rescheduleErr := s.repository.Reschedule(ctx, event.EventID, s.now().Add(s.backoff), err.Error()); rescheduleErr != nil {
				publishErr = errors.Join(publishErr, err, rescheduleErr)
				continue
			}

			publishErr = errors.Join(publishErr, err)
			continue
		}

		publishedAt := s.now()
		if err := s.repository.MarkPublished(ctx, event.EventID, publishedAt); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, EventRepositoryPublish)
			return err
		}

		s.logger.InfowCtx(ctx, LogOutboxEventPublished,
			LogKeyEventID, event.EventID,
			LogKeyEventType, event.EventType,
			LogKeyAggregateType, event.AggregateType,
			LogKeyAggregateID, event.AggregateID,
		)
	}

	if publishErr != nil {
		span.RecordError(publishErr)
		span.SetStatus(codes.Error, EventPublish)
		return publishErr
	}

	span.SetStatus(codes.Ok, EventSuccess)
	return nil
}
