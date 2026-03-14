package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/realtime/core/domain"
)

type projectionReadyEnvelope struct {
	EventType       string `json:"event_type"`
	EventVersion    string `json:"event_version"`
	UserID          uint64 `json:"user_id"`
	RecordID        uint64 `json:"record_id"`
	SourceEventID   string `json:"source_event_id"`
	SourceEventType string `json:"source_event_type"`
	ProjectedAtUTC  string `json:"projected_at_utc"`
	TraceID         string `json:"trace_id"`
	RequestID       string `json:"request_id"`
}

func (r *ProjectionEventReader) Read(ctx context.Context) (domain.Event, error) {
	message, err := r.reader.ReadMessage(ctx)
	if err != nil {
		return domain.Event{}, err
	}

	return parseProjectionReadyPayload(message.Value)
}

func parseProjectionReadyPayload(payload []byte) (domain.Event, error) {
	var envelope projectionReadyEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return domain.Event{}, fmt.Errorf("decode projection ready event: %w", err)
	}

	projectedAtUTC, err := time.Parse(time.RFC3339Nano, envelope.ProjectedAtUTC)
	if err != nil {
		return domain.Event{}, fmt.Errorf("parse projected_at_utc: %w", err)
	}

	return domain.Event{
		Type:           "record_projection_changed",
		UserID:         envelope.UserID,
		RecordID:       envelope.RecordID,
		Action:         actionFromEventType(envelope.EventType),
		ProjectedAtUTC: projectedAtUTC.UTC(),
		SourceEventID:  envelope.SourceEventID,
		TraceID:        envelope.TraceID,
		RequestID:      envelope.RequestID,
	}, nil
}

func actionFromEventType(eventType string) string {
	switch eventType {
	case "record.projection.created":
		return "created"
	case "record.projection.deleted":
		return "deleted"
	default:
		return "updated"
	}
}
