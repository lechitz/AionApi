package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/lechitz/AionApi/internal/eventoutbox/core/domain"
)

func (p *EventPublisher) topicFor(event domain.Event) (string, error) {
	switch event.AggregateType {
	case RecordAggregateType:
		return p.recordEventsTopic, nil
	default:
		return "", fmt.Errorf("%s: %s", ErrUnsupportedAggregateType, event.AggregateType)
	}
}

func buildEnvelope(event domain.Event) envelope {
	return envelope{
		EventID:       event.EventID,
		AggregateType: event.AggregateType,
		AggregateID:   event.AggregateID,
		EventType:     event.EventType,
		EventVersion:  event.EventVersion,
		Source:        event.Source,
		TraceID:       event.TraceID,
		RequestID:     event.RequestID,
		OccurredAtUTC: event.CreatedAt.UTC().Format("2006-01-02T15:04:05.000000Z07:00"),
		Payload:       json.RawMessage(event.PayloadJSON),
	}
}
