// Package kafka publishes canonical outbox events to Kafka topics.
package kafka

const (
	// PublisherTracerName is the tracer name used by the Kafka outbox publisher adapter.
	PublisherTracerName = "aionapi.eventoutbox.kafka.publisher"
)

const (
	// SpanOutboxKafkaPublish is the span name for Kafka publication.
	SpanOutboxKafkaPublish = "eventoutbox.kafka.publish"
)

const (
	// OpOutboxKafkaPublish is the operation name for Kafka publication.
	OpOutboxKafkaPublish = "event_outbox_kafka_publish"
)

const (
	// RecordAggregateType identifies record events emitted by the canonical API.
	RecordAggregateType = "record"
)

const (
	// LogPublishingEvent indicates one outbox event is being sent to Kafka.
	LogPublishingEvent = "publishing outbox event to kafka"
	// LogPublishedEvent indicates one outbox event was sent to Kafka.
	LogPublishedEvent = "outbox event published to kafka"
)

const (
	// ErrUnsupportedAggregateType is returned when no topic is configured for the aggregate.
	ErrUnsupportedAggregateType = "unsupported aggregate type for kafka publication"
)
