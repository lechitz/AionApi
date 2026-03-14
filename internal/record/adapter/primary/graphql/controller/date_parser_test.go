package controller_test

import (
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/adapter/primary/graphql/controller"
	"github.com/stretchr/testify/require"
)

func TestParseRecordsDayQueryWithLocation_AbsoluteDates(t *testing.T) {
	loc := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, 2, 20, 10, 0, 0, 0, loc)

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "yyyy-mm-dd",
			input:    "2026-02-14",
			expected: time.Date(2026, 2, 14, 0, 0, 0, 0, loc),
		},
		{
			name:     "dd/mm/yyyy",
			input:    "14/02/2026",
			expected: time.Date(2026, 2, 14, 0, 0, 0, 0, loc),
		},
		{
			name:     "rfc3339",
			input:    "2026-02-14T13:20:00Z",
			expected: time.Date(2026, 2, 14, 0, 0, 0, 0, loc),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := controller.ParseRecordsDayQueryWithLocation(tc.input, now, loc)
			require.NoError(t, err)
			require.True(t, got.Equal(tc.expected), "got=%s expected=%s", got, tc.expected)
		})
	}
}

func TestParseRecordsDayQueryWithLocation_NaturalLanguage(t *testing.T) {
	loc := time.FixedZone("BRT", -3*60*60)
	nowMonday := time.Date(2026, 2, 16, 10, 0, 0, 0, loc) // Monday
	nowFriday := time.Date(2026, 2, 20, 10, 0, 0, 0, loc) // Friday

	tests := []struct {
		name     string
		input    string
		now      time.Time
		expected time.Time
	}{
		{
			name:     "empty means today",
			input:    "",
			now:      nowMonday,
			expected: time.Date(2026, 2, 16, 0, 0, 0, 0, loc),
		},
		{
			name:     "hoje",
			input:    "hoje",
			now:      nowMonday,
			expected: time.Date(2026, 2, 16, 0, 0, 0, 0, loc),
		},
		{
			name:     "ontem",
			input:    "ontem",
			now:      nowMonday,
			expected: time.Date(2026, 2, 15, 0, 0, 0, 0, loc),
		},
		{
			name:     "weekday defaults to most recent occurrence",
			input:    "sexta",
			now:      nowMonday,
			expected: time.Date(2026, 2, 13, 0, 0, 0, 0, loc),
		},
		{
			name:     "sexta passada from monday",
			input:    "sexta-feira passada",
			now:      nowMonday,
			expected: time.Date(2026, 2, 13, 0, 0, 0, 0, loc),
		},
		{
			name:     "sexta on friday means today",
			input:    "sexta",
			now:      nowFriday,
			expected: time.Date(2026, 2, 20, 0, 0, 0, 0, loc),
		},
		{
			name:     "sexta passada on friday means previous week",
			input:    "sexta passada",
			now:      nowFriday,
			expected: time.Date(2026, 2, 13, 0, 0, 0, 0, loc),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := controller.ParseRecordsDayQueryWithLocation(tc.input, tc.now, loc)
			require.NoError(t, err)
			require.True(t, got.Equal(tc.expected), "got=%s expected=%s", got, tc.expected)
		})
	}
}
