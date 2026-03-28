package kafka

import (
	"strings"

	"github.com/lechitz/aion-api/internal/platform/config"
	"github.com/lechitz/aion-api/internal/platform/ports/output/logger"
	kafkago "github.com/segmentio/kafka-go"
)

// EventPublisher publishes durable outbox events to Kafka.
type EventPublisher struct {
	writer            *kafkago.Writer
	logger            logger.ContextLogger
	recordEventsTopic string
}

// NewEventPublisher creates a Kafka-backed outbox publisher.
func NewEventPublisher(cfg config.KafkaConfig, log logger.ContextLogger) *EventPublisher {
	return &EventPublisher{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(strings.Split(cfg.Brokers, ",")...),
			Balancer:     &kafkago.LeastBytes{},
			RequiredAcks: kafkago.RequireAll,
			Async:        false,
		},
		logger:            log,
		recordEventsTopic: cfg.RecordEventsTopic,
	}
}

// Close releases the underlying Kafka writer.
func (p *EventPublisher) Close() error {
	return p.writer.Close()
}
