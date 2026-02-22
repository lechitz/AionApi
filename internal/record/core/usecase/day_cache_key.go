package usecase

import "time"

// cacheDayStart normalizes an instant to UTC midnight so cache keys stay stable.
func cacheDayStart(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}
