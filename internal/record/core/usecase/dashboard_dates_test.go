package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNormalizeDashboardDate_PreservesExplicitCalendarDate(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	explicitDate, err := time.Parse("2006-01-02", "2026-03-18")
	require.NoError(t, err)

	got := normalizeDashboardDate(explicitDate, loc)

	require.Equal(t, time.Date(2026, 3, 18, 0, 0, 0, 0, loc), got)
}

func TestNormalizeInsightDate_PreservesExplicitCalendarDate(t *testing.T) {
	t.Parallel()

	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	explicitDate, err := time.Parse("2006-01-02", "2026-03-18")
	require.NoError(t, err)

	got := normalizeInsightDate(explicitDate, loc)

	require.Equal(t, time.Date(2026, 3, 18, 0, 0, 0, 0, loc), got)
}
