// Package domain contains domain entities for the event outbox bounded context.
package domain

import "time"

// Event represents one canonical domain event persisted in the outbox.
type Event struct {
	EventID        string
	AggregateType  string
	AggregateID    string
	EventType      string
	EventVersion   string
	Source         string
	Status         string
	TraceID        string
	RequestID      string
	PayloadJSON    []byte
	AttemptCount   int
	AvailableAtUTC time.Time
	PublishedAtUTC *time.Time
	LastError      string
	CreatedAt      time.Time
}
