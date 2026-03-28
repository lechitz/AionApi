package usecase_test

import (
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_DashboardSnapshot_BuildsChecklistWithTarget(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	targetDate := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	startUTC := time.Date(2026, 3, 18, 0, 0, 0, 0, loc).UTC()
	endUTC := time.Date(2026, 3, 18, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()
	goal := 3.0

	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{{
			ID:          8,
			UserID:      userID,
			MetricKey:   "intentions",
			DisplayName: "Intencoes",
			TagID:       8,
			TagIDs:      []uint64{8},
			ValueSource: usecase.DashboardValueSourceCount,
			Aggregation: usecase.DashboardAggregationSum,
			Unit:        usecase.DashboardUnitCount,
			GoalDefault: &goal,
			IsActive:    true,
		}}, nil)
	suite.RecordRepository.EXPECT().
		ListAllBetween(gomock.Any(), userID, startUTC, endUTC, gomock.Any()).
		Return([]domain.Record{
			{ID: 1, UserID: userID, TagID: 8, EventTime: time.Date(2026, 3, 18, 11, 0, 0, 0, time.UTC)},
			{ID: 2, UserID: userID, TagID: 8, EventTime: time.Date(2026, 3, 18, 13, 0, 0, 0, time.UTC)},
		}, nil)
	suite.RecordRepository.EXPECT().
		ListGoalTemplates(gomock.Any(), userID).
		Return([]domain.GoalTemplate{}, nil)

	out, err := suite.RecordService.DashboardSnapshot(t.Context(), userID, input.DashboardSnapshotQuery{
		Date:     targetDate,
		Timezone: "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, out.Metrics, 1)
	require.NotNil(t, out.Metrics[0].Checklist)
	assert.Equal(t, domain.DashboardChecklistModeCountGoal, out.Metrics[0].Checklist.Mode)
	assert.Equal(t, 2, out.Metrics[0].Checklist.CompletedCount)
	require.NotNil(t, out.Metrics[0].Checklist.TargetCount)
	assert.Equal(t, 3, *out.Metrics[0].Checklist.TargetCount)
	require.NotNil(t, out.Metrics[0].Checklist.RemainingCount)
	assert.Equal(t, 1, *out.Metrics[0].Checklist.RemainingCount)
	assert.InDelta(t, 2.0/3.0, out.Metrics[0].Checklist.CompletionRatio, 0.000001)
}

func TestService_DashboardSnapshot_BuildsChecklistWithoutTarget(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	targetDate := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	startUTC := time.Date(2026, 3, 18, 0, 0, 0, 0, loc).UTC()
	endUTC := time.Date(2026, 3, 18, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()

	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{{
			ID:          9,
			UserID:      userID,
			MetricKey:   "cafe_preto_count",
			DisplayName: "Cafe Preto",
			TagID:       25,
			TagIDs:      []uint64{25},
			ValueSource: usecase.DashboardValueSourceCount,
			Aggregation: usecase.DashboardAggregationSum,
			Unit:        usecase.DashboardUnitCount,
			IsActive:    true,
		}}, nil)
	suite.RecordRepository.EXPECT().
		ListAllBetween(gomock.Any(), userID, startUTC, endUTC, gomock.Any()).
		Return([]domain.Record{
			{ID: 1, UserID: userID, TagID: 25, EventTime: time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)},
			{ID: 2, UserID: userID, TagID: 25, EventTime: time.Date(2026, 3, 18, 15, 0, 0, 0, time.UTC)},
			{ID: 3, UserID: userID, TagID: 25, EventTime: time.Date(2026, 3, 18, 18, 0, 0, 0, time.UTC)},
		}, nil)
	suite.RecordRepository.EXPECT().
		ListGoalTemplates(gomock.Any(), userID).
		Return([]domain.GoalTemplate{}, nil)

	out, err := suite.RecordService.DashboardSnapshot(t.Context(), userID, input.DashboardSnapshotQuery{
		Date:     targetDate,
		Timezone: "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, out.Metrics, 1)
	require.NotNil(t, out.Metrics[0].Checklist)
	assert.Equal(t, domain.DashboardChecklistModeCountOnly, out.Metrics[0].Checklist.Mode)
	assert.Equal(t, 3, out.Metrics[0].Checklist.CompletedCount)
	assert.Nil(t, out.Metrics[0].Checklist.TargetCount)
	assert.Nil(t, out.Metrics[0].Checklist.RemainingCount)
	assert.InDelta(t, 1.0, out.Metrics[0].Checklist.CompletionRatio, 0.000001)
}

func TestService_DashboardSnapshot_BuildsZeroChecklistState(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	targetDate := time.Date(2026, 3, 18, 12, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	startUTC := time.Date(2026, 3, 18, 0, 0, 0, 0, loc).UTC()
	endUTC := time.Date(2026, 3, 18, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()
	goal := 3.0

	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{{
			ID:          8,
			UserID:      userID,
			MetricKey:   "intentions",
			DisplayName: "Intencoes",
			TagID:       8,
			TagIDs:      []uint64{8},
			ValueSource: usecase.DashboardValueSourceCount,
			Aggregation: usecase.DashboardAggregationSum,
			Unit:        usecase.DashboardUnitCount,
			GoalDefault: &goal,
			IsActive:    true,
		}}, nil)
	suite.RecordRepository.EXPECT().
		ListAllBetween(gomock.Any(), userID, startUTC, endUTC, gomock.Any()).
		Return([]domain.Record{}, nil)
	suite.RecordRepository.EXPECT().
		ListGoalTemplates(gomock.Any(), userID).
		Return([]domain.GoalTemplate{}, nil)

	out, err := suite.RecordService.DashboardSnapshot(t.Context(), userID, input.DashboardSnapshotQuery{
		Date:     targetDate,
		Timezone: "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, out.Metrics, 1)
	require.NotNil(t, out.Metrics[0].Checklist)
	assert.Equal(t, 0, out.Metrics[0].Checklist.CompletedCount)
	require.NotNil(t, out.Metrics[0].Checklist.RemainingCount)
	assert.Equal(t, 3, *out.Metrics[0].Checklist.RemainingCount)
	assert.InDelta(t, 0.0, out.Metrics[0].Checklist.CompletionRatio, 0.000001)
}

func TestService_DashboardSnapshot_RespectsTimezoneBoundaryForChecklist(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	targetDate := time.Date(2026, 3, 10, 15, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	startUTC := time.Date(2026, 3, 10, 0, 0, 0, 0, loc).UTC()
	endUTC := time.Date(2026, 3, 10, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()
	goal := 2.0

	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{{
			ID:          8,
			UserID:      userID,
			MetricKey:   "intentions",
			DisplayName: "Intencoes",
			TagID:       8,
			TagIDs:      []uint64{8},
			ValueSource: usecase.DashboardValueSourceCount,
			Aggregation: usecase.DashboardAggregationSum,
			Unit:        usecase.DashboardUnitCount,
			GoalDefault: &goal,
			IsActive:    true,
		}}, nil)
	suite.RecordRepository.EXPECT().
		ListAllBetween(gomock.Any(), userID, startUTC, endUTC, gomock.Any()).
		Return([]domain.Record{
			{ID: 2, UserID: userID, TagID: 8, EventTime: time.Date(2026, 3, 10, 3, 30, 0, 0, time.UTC)}, // 2026-03-10 00:30 BRT
		}, nil)
	suite.RecordRepository.EXPECT().
		ListGoalTemplates(gomock.Any(), userID).
		Return([]domain.GoalTemplate{}, nil)

	out, err := suite.RecordService.DashboardSnapshot(t.Context(), userID, input.DashboardSnapshotQuery{
		Date:     targetDate,
		Timezone: "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, out.Metrics, 1)
	require.NotNil(t, out.Metrics[0].Checklist)
	assert.Equal(t, 1, out.Metrics[0].Checklist.CompletedCount)
}
