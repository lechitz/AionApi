package controller

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDateOrDefault_EmptyReturnsZeroTime(t *testing.T) {
	t.Parallel()

	got, err := parseDateOrDefault("")
	require.NoError(t, err)
	require.True(t, got.IsZero())
}

func TestParseDateOrDefault_ParsesCalendarDate(t *testing.T) {
	t.Parallel()

	got, err := parseDateOrDefault("2026-03-18")
	require.NoError(t, err)
	require.Equal(t, time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC), got)
}
