package usecase_test

import (
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/ports/input"
	tagdomain "github.com/lechitz/AionApi/internal/tag/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_InsightFeed_EmitsInsufficientDataEvidence(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	date := time.Date(2026, 3, 10, 15, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	currentStart := time.Date(2026, 3, 4, 0, 0, 0, 0, loc).UTC()
	currentEnd := time.Date(2026, 3, 10, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()
	prevStart := time.Date(2026, 2, 25, 0, 0, 0, 0, loc).UTC()
	prevEnd := currentStart.Add(-time.Nanosecond)

	gomock.InOrder(
		suite.RecordRepository.EXPECT().
			ListAllBetween(gomock.Any(), userID, currentStart, currentEnd, gomock.Any()).
			Return([]domain.Record{}, nil),
		suite.RecordRepository.EXPECT().
			ListAllBetween(gomock.Any(), userID, prevStart, prevEnd, gomock.Any()).
			Return([]domain.Record{
				{ID: 1, UserID: userID, TagID: 10, EventTime: time.Date(2026, 2, 27, 12, 0, 0, 0, time.UTC)},
				{ID: 2, UserID: userID, TagID: 10, EventTime: time.Date(2026, 2, 28, 12, 0, 0, 0, time.UTC)},
				{ID: 3, UserID: userID, TagID: 10, EventTime: time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)},
			}, nil),
	)

	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{{TagID: 10, TagIDs: []uint64{10}, DisplayName: "Agua"}}, nil)
	suite.TagRepository.EXPECT().
		GetAll(gomock.Any(), userID).
		Return([]tagdomain.Tag{{ID: 10, UserID: userID, CategoryID: 1, Name: "Agua"}}, nil)

	got, err := suite.RecordService.InsightFeed(t.Context(), userID, input.InsightFeedQuery{
		Window:   string(domain.InsightWindow7D),
		Limit:    5,
		Date:     date,
		Timezone: "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.Equal(t, "activity_gap", got[0].Type)
	assert.Equal(t, "Sem atividade na janela analisada", got[0].Title)
	require.Len(t, got[0].Evidence, 2)
	assert.Equal(t, "janela final", got[0].Evidence[0].Label)
	assert.Equal(t, "2026-03-10", got[0].Evidence[0].Value)
	assert.Equal(t, "registros", got[0].Evidence[1].Label)
	assert.Equal(t, "0", got[0].Evidence[1].Value)

	assert.Equal(t, "recent_change", got[1].Type)
	assert.Equal(t, "Ritmo abaixo da janela anterior", got[1].Title)
	require.Len(t, got[1].Evidence, 2)
	assert.Equal(t, "janela atual", got[1].Evidence[0].Label)
	assert.Equal(t, "0", got[1].Evidence[0].Value)
	assert.Equal(t, "janela anterior", got[1].Evidence[1].Label)
	assert.Equal(t, "3", got[1].Evidence[1].Value)
}

func TestService_AnalyticsSeries_RespectsTimezoneBoundary(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	date := time.Date(2026, 3, 10, 15, 0, 0, 0, time.UTC)
	loc, err := time.LoadLocation("America/Sao_Paulo")
	require.NoError(t, err)

	currentStart := time.Date(2026, 3, 4, 0, 0, 0, 0, loc).UTC()
	currentEnd := time.Date(2026, 3, 10, 23, 59, 59, int(time.Second-time.Nanosecond), loc).UTC()

	suite.RecordRepository.EXPECT().
		ListAllBetween(gomock.Any(), userID, currentStart, currentEnd, gomock.Any()).
		Return([]domain.Record{
			{ID: 1, UserID: userID, TagID: 10, EventTime: time.Date(2026, 3, 10, 2, 30, 0, 0, time.UTC)}, // 2026-03-09 23:30 BRT
			{ID: 2, UserID: userID, TagID: 10, EventTime: time.Date(2026, 3, 10, 3, 30, 0, 0, time.UTC)}, // 2026-03-10 00:30 BRT
		}, nil)
	suite.RecordRepository.EXPECT().
		ListMetricDefinitions(gomock.Any(), userID).
		Return([]domain.MetricDefinition{}, nil)
	suite.TagRepository.EXPECT().
		GetAll(gomock.Any(), userID).
		Return([]tagdomain.Tag{{ID: 10, UserID: userID, CategoryID: 1, Name: "Agua"}}, nil)

	got, err := suite.RecordService.AnalyticsSeries(t.Context(), userID, input.AnalyticsSeriesQuery{
		SeriesKey: "records.count",
		Window:    string(domain.InsightWindow7D),
		Date:      date,
		Timezone:  "America/Sao_Paulo",
	})

	require.NoError(t, err)
	require.Len(t, got.Points, 7)
	require.NotNil(t, got.Summary)
	assert.Equal(t, "records.count across 7 days", *got.Summary)
	require.NotNil(t, got.Points[5].Value)
	require.NotNil(t, got.Points[6].Value)
	assert.InDelta(t, 1.0, *got.Points[5].Value, 0.000001)
	assert.Equal(t, "2026-03-09", *got.Points[5].Label)
	assert.InDelta(t, 1.0, *got.Points[6].Value, 0.000001)
	assert.Equal(t, "2026-03-10", *got.Points[6].Label)
}
