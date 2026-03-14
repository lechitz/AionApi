//nolint:testpackage // Tests validate package-private parsing helpers directly.
package kafka

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/realtime/core/domain"
)

func TestActionFromEventType(t *testing.T) {
	tests := []struct {
		eventType string
		want      string
	}{
		{eventType: "record.projection.created", want: "created"},
		{eventType: "record.projection.updated", want: "updated"},
		{eventType: "record.projection.deleted", want: "deleted"},
	}

	for _, tt := range tests {
		if got := actionFromEventType(tt.eventType); got != tt.want {
			t.Fatalf("eventType=%s want=%s got=%s", tt.eventType, tt.want, got)
		}
	}
}

func TestParseProjectionReadyPayload(t *testing.T) {
	projectedAt := time.Now().UTC().Truncate(time.Second)
	payload, err := json.Marshal(projectionReadyEnvelope{
		EventType:       "record.projection.created",
		EventVersion:    "v1",
		UserID:          14,
		RecordID:        42,
		SourceEventID:   "evt-1",
		SourceEventType: "record.created",
		ProjectedAtUTC:  projectedAt.Format(time.RFC3339Nano),
		TraceID:         "trace-1",
		RequestID:       "req-1",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	got, err := parseProjectionReadyPayload(payload)
	if err != nil {
		t.Fatalf("parse payload: %v", err)
	}

	want := domain.Event{
		Type:           "record_projection_changed",
		UserID:         14,
		RecordID:       42,
		Action:         "created",
		ProjectedAtUTC: projectedAt,
		SourceEventID:  "evt-1",
		TraceID:        "trace-1",
		RequestID:      "req-1",
	}

	if got != want {
		t.Fatalf("unexpected event: %#v", got)
	}
}
