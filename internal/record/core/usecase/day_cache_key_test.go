package usecase_test

import (
	"testing"
	"time"

	usecase "github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/stretchr/testify/require"
)

func TestCacheDayStart_NormalizesToUTCMidnight(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "sao paulo morning stays same UTC date",
			input:    time.Date(2026, 2, 22, 10, 20, 0, 0, time.FixedZone("-03", -3*3600)),
			expected: time.Date(2026, 2, 22, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "positive offset can move to previous UTC date",
			input:    time.Date(2026, 2, 22, 0, 30, 0, 0, time.FixedZone("+03", 3*3600)),
			expected: time.Date(2026, 2, 21, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := usecase.CacheDayStart(tc.input)
			require.Equal(t, tc.expected, got)
			require.Equal(t, time.UTC, got.Location())
		})
	}
}
