package kafka

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/eventoutbox/core/domain"
)

func TestTopicForRecordEvent(t *testing.T) {
	t.Parallel()

	publisher := &EventPublisher{recordEventsTopic: "aion.record.events.v1"}
	topic, err := publisher.topicFor(domain.Event{AggregateType: RecordAggregateType})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if topic != "aion.record.events.v1" {
		t.Fatalf("expected record topic, got %q", topic)
	}
}

func TestBuildEnvelope(t *testing.T) {
	t.Parallel()

	event := domain.Event{
		EventID:       "evt-1",
		AggregateType: "record",
		AggregateID:   "123",
		EventType:     "record.created",
		EventVersion:  "v1",
		Source:        "aion-api",
		TraceID:       "trace-1",
		RequestID:     "req-1",
		CreatedAt:     time.Date(2026, time.March, 13, 13, 0, 0, 0, time.UTC),
		PayloadJSON:   []byte(`{"record_id":123}`),
	}

	envelope := buildEnvelope(event)
	if envelope.EventID != event.EventID {
		t.Fatalf("expected event id %q, got %q", event.EventID, envelope.EventID)
	}
	if string(envelope.Payload) != `{"record_id":123}` {
		t.Fatalf("expected payload passthrough, got %s", string(envelope.Payload))
	}
	if envelope.OccurredAtUTC != "2026-03-13T13:00:00.000000Z" {
		t.Fatalf("unexpected occurred_at_utc: %q", envelope.OccurredAtUTC)
	}
}
