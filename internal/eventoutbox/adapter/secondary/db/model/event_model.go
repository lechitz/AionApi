// Package model contains database models for the event outbox context.
package model

import "time"

const (
	// EventOutboxTable is the fully qualified table name for durable outbox events.
	EventOutboxTable = "aion_api.event_outbox"
)

// EventDB represents one durable outbox event row.
type EventDB struct {
	ID             uint64     `gorm:"primaryKey;column:id;autoIncrement"`
	EventID        string     `gorm:"column:event_id;type:uuid;not null;uniqueIndex"`
	AggregateType  string     `gorm:"column:aggregate_type;size:64;not null;index"`
	AggregateID    string     `gorm:"column:aggregate_id;size:128;not null;index"`
	EventType      string     `gorm:"column:event_type;size:128;not null"`
	EventVersion   string     `gorm:"column:event_version;size:16;not null"`
	Source         string     `gorm:"column:source;size:32;not null"`
	TraceID        string     `gorm:"column:trace_id;size:64;index"`
	RequestID      string     `gorm:"column:request_id;size:64"`
	Status         string     `gorm:"column:status;size:24;not null;index"`
	AttemptCount   int        `gorm:"column:attempt_count;not null;default:0"`
	AvailableAtUTC time.Time  `gorm:"column:available_at_utc;not null;index"`
	PublishedAtUTC *time.Time `gorm:"column:published_at_utc"`
	LastError      string     `gorm:"column:last_error"`
	PayloadJSON    []byte     `gorm:"column:payload_json;type:jsonb;not null"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName returns the fully qualified database table name for EventDB.
func (EventDB) TableName() string {
	return EventOutboxTable
}
