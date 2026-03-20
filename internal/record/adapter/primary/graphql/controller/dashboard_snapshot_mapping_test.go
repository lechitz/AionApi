package controller_test

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	"github.com/stretchr/testify/require"
)

func TestDashboardSnapshot_MapsCanonicalChecklistPayload(t *testing.T) {
	t.Parallel()

	h, ctrl := newRecordController(t, &recordServiceStub{
		dashboardFn: func(_ context.Context, userID uint64, query input.DashboardSnapshotQuery) (domain.DashboardSnapshot, error) {
			require.Equal(t, uint64(999), userID)
			require.Equal(t, "America/Sao_Paulo", query.Timezone)
			target := 3
			remaining := 1
			return domain.DashboardSnapshot{
				Date:     time.Date(2026, 3, 18, 0, 0, 0, 0, time.UTC),
				Timezone: "America/Sao_Paulo",
				Metrics: []domain.DashboardMetricValue{{
					MetricKey:   "intentions",
					Label:       "Intencoes",
					Value:       2,
					Unit:        "count",
					ProgressPct: 66.67,
					Status:      "pending",
					Checklist: &domain.DashboardChecklistValue{
						MetricKey:       "intentions",
						Label:           "Intencoes",
						CompletedCount:  2,
						TargetCount:     &target,
						RemainingCount:  &remaining,
						CompletionRatio: 2.0 / 3.0,
						Status:          "pending",
						Mode:            domain.DashboardChecklistModeCountGoal,
					},
				}},
			}, nil
		},
	})
	defer ctrl.Finish()

	tz := "America/Sao_Paulo"
	out, err := h.DashboardSnapshot(t.Context(), 999, "2026-03-18", &tz)
	require.NoError(t, err)
	require.Len(t, out.Metrics, 1)
	require.NotNil(t, out.Metrics[0].Checklist)
	require.EqualValues(t, 2, out.Metrics[0].Checklist.CompletedCount)
	require.NotNil(t, out.Metrics[0].Checklist.TargetCount)
	require.EqualValues(t, 3, *out.Metrics[0].Checklist.TargetCount)
	require.Equal(t, domain.DashboardChecklistModeCountGoal, out.Metrics[0].Checklist.Mode)
}
