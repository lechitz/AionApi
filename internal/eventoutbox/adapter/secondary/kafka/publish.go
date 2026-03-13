package kafka

import (
	"context"
	"encoding/json"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Publish sends one outbox event to the configured Kafka topic.
func (p *EventPublisher) Publish(ctx context.Context, event domain.Event) error {
	topic, err := p.topicFor(event)
	if err != nil {
		return err
	}

	tr := otel.Tracer(PublisherTracerName)
	ctx, span := tr.Start(ctx, SpanOutboxKafkaPublish, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OpOutboxKafkaPublish),
		attribute.String(commonkeys.Entity, event.AggregateType),
		attribute.String("aggregate_id", event.AggregateID),
		attribute.String("event_type", event.EventType),
		attribute.String("event_id", event.EventID),
		attribute.String("kafka_topic", topic),
	))
	defer span.End()

	payload, err := json.Marshal(buildEnvelope(event))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxKafkaPublish)
		return err
	}

	msg := kafkago.Message{
		Topic: topic,
		Key:   []byte(event.AggregateID),
		Value: payload,
		Headers: []kafkago.Header{
			{Key: "event_id", Value: []byte(event.EventID)},
			{Key: "event_type", Value: []byte(event.EventType)},
			{Key: "event_version", Value: []byte(event.EventVersion)},
			{Key: "aggregate_type", Value: []byte(event.AggregateType)},
			{Key: "aggregate_id", Value: []byte(event.AggregateID)},
			{Key: "source", Value: []byte(event.Source)},
		},
	}

	p.logger.InfowCtx(ctx, LogPublishingEvent,
		"event_id", event.EventID,
		"aggregate_type", event.AggregateType,
		"aggregate_id", event.AggregateID,
		"event_type", event.EventType,
		"kafka_topic", topic,
	)

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, OpOutboxKafkaPublish)
		return err
	}

	span.SetStatus(codes.Ok, LogPublishedEvent)
	p.logger.InfowCtx(ctx, LogPublishedEvent,
		"event_id", event.EventID,
		"aggregate_type", event.AggregateType,
		"aggregate_id", event.AggregateID,
		"event_type", event.EventType,
		"kafka_topic", topic,
	)

	return nil
}
