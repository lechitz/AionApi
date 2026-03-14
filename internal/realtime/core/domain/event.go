// Package domain defines the realtime event model exposed by the API.
package domain

import "time"

// Event is the payload delivered to realtime subscribers.
type Event struct {
	Type           string    `json:"type"`
	UserID         uint64    `json:"userId"`
	RecordID       uint64    `json:"recordId"`
	Action         string    `json:"action"`
	ProjectedAtUTC time.Time `json:"projectedAtUTC"`
	SourceEventID  string    `json:"sourceEventId,omitempty"`
	TraceID        string    `json:"traceId,omitempty"`
	RequestID      string    `json:"requestId,omitempty"`
}
