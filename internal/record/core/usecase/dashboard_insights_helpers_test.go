package usecase

import (
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/stretchr/testify/require"
)

func TestResolveInsightLocation_Fallbacks(t *testing.T) {
	t.Run("uses default timezone when empty", func(t *testing.T) {
		loc, tz := resolveInsightLocation("")
		require.Equal(t, DefaultTimezone, tz)
		require.Equal(t, DefaultTimezone, loc.String())
	})

	t.Run("falls back to UTC when timezone is invalid", func(t *testing.T) {
		loc, tz := resolveInsightLocation("Invalid/Timezone")
		require.Equal(t, "UTC", tz)
		require.Equal(t, time.UTC, loc)
	})
}

func TestBuildActivityGapInsight_IncludesEvidence(t *testing.T) {
	loc := time.FixedZone("BRT", -3*60*60)
	now := time.Date(2026, 3, 11, 1, 0, 0, 0, time.UTC)
	targetDate := time.Date(2026, 3, 10, 0, 0, 0, 0, loc)

	got := buildActivityGapInsight(nil, domain.InsightWindow7D, now, loc, targetDate)
	require.NotNil(t, got)
	require.Equal(t, "activity_gap", got.Type)
	require.Equal(t, "Sem atividade na janela analisada", got.Title)
	require.Len(t, got.Evidence, 2)
	require.Equal(t, "janela final", got.Evidence[0].Label)
	require.Equal(t, "2026-03-10", got.Evidence[0].Value)
	require.Equal(t, "date", got.Evidence[0].Kind)
	require.Equal(t, "registros", got.Evidence[1].Label)
	require.Equal(t, "0", got.Evidence[1].Value)
	require.Equal(t, "count", got.Evidence[1].Kind)
}

func TestBuildRecentChangeInsight_IncludesEvidence(t *testing.T) {
	now := time.Date(2026, 3, 11, 1, 0, 0, 0, time.UTC)
	currRecords := []domain.Record{
		{ID: 1, EventTime: now},
	}
	prevRecords := []domain.Record{
		{ID: 2, EventTime: now.AddDate(0, 0, -8)},
		{ID: 3, EventTime: now.AddDate(0, 0, -9)},
		{ID: 4, EventTime: now.AddDate(0, 0, -10)},
		{ID: 5, EventTime: now.AddDate(0, 0, -11)},
	}

	got := buildRecentChangeInsight(currRecords, prevRecords, domain.InsightWindow7D, now, 7)
	require.NotNil(t, got)
	require.Equal(t, "recent_change", got.Type)
	require.Equal(t, "Ritmo abaixo da janela anterior", got.Title)
	require.Len(t, got.Evidence, 2)
	require.Equal(t, "janela atual", got.Evidence[0].Label)
	require.Equal(t, "1", got.Evidence[0].Value)
	require.Equal(t, "count", got.Evidence[0].Kind)
	require.Equal(t, "janela anterior", got.Evidence[1].Label)
	require.Equal(t, "4", got.Evidence[1].Value)
	require.Equal(t, "count", got.Evidence[1].Kind)
}
