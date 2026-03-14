package domain

import "time"

// Stats aggregates operational counters and timestamps for the event outbox.
type Stats struct {
	PendingCount       int64
	PublishedCount     int64
	FailedCount        int64
	OldestPendingAtUTC *time.Time
}
