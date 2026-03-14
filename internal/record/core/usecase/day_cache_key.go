package usecase

import "time"

// CacheDayStart normalizes an instant to UTC midnight so cache keys stay stable.
func CacheDayStart(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
