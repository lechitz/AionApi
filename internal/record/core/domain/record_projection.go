package domain

import "time"

// RecordProjection is the canonical read model materialized from the derived Kafka projection.
type RecordProjection struct {
	RecordID           uint64
	UserID             uint64
	TagID              uint64
	Description        *string
	EventTimeUTC       time.Time
	RecordedAtUTC      *time.Time
	Status             *string
	Timezone           *string
	DurationSeconds    *int
	Value              *float64
	Source             *string
	LastEventID        string
	LastEventType      string
	LastEventVersion   string
	LastTraceID        *string
	LastRequestID      *string
	LastKafkaTopic     string
	LastKafkaPartition int
	LastKafkaOffset    int64
	LastConsumedAtUTC  time.Time
	PayloadJSON        []byte
	CreatedAtUTC       time.Time
	UpdatedAtUTC       time.Time
}
