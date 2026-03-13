package kafka

import "encoding/json"

type envelope struct {
	EventID       string          `json:"event_id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	EventVersion  string          `json:"event_version"`
	Source        string          `json:"source"`
	TraceID       string          `json:"trace_id,omitempty"`
	RequestID     string          `json:"request_id,omitempty"`
	OccurredAtUTC string          `json:"occurred_at_utc"`
	Payload       json.RawMessage `json:"payload"`
}
