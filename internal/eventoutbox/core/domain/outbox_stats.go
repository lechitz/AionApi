package domain

import "time"

type Stats struct {
	PendingCount       int64
	PublishedCount     int64
	FailedCount        int64
	OldestPendingAtUTC *time.Time
}
