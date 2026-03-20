package usecase

import "time"

func normalizeDashboardDate(date time.Time, loc *time.Location) time.Time {
	if date.IsZero() {
		date = time.Now().In(loc)
	}

	// GraphQL date inputs arrive as YYYY-MM-DD and are parsed as midnight UTC.
	// Treat those values as calendar dates in the requested location instead of
	// shifting the day boundary through an instant conversion.
	if date.Location() == time.UTC &&
		date.Hour() == 0 &&
		date.Minute() == 0 &&
		date.Second() == 0 &&
		date.Nanosecond() == 0 {
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)
	}

	local := date.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
}
