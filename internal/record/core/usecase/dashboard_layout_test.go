package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/ports/input"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_UpsertDashboardWidget_AllowsInsightFeedWithoutMetricDefinition(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	viewID := uint64(10)

	suite.RecordRepository.EXPECT().
		ListDashboardWidgetsByView(gomock.Any(), userID, viewID).
		Return([]domain.DashboardWidget{
			{ID: 1, UserID: userID, ViewID: viewID},
			{ID: 2, UserID: userID, ViewID: viewID},
		}, nil)

	suite.RecordRepository.EXPECT().
		UpsertDashboardWidget(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, widget domain.DashboardWidget) (domain.DashboardWidget, error) {
			assert.Equal(t, userID, widget.UserID)
			assert.Equal(t, viewID, widget.ViewID)
			assert.Equal(t, domain.DashboardWidgetTypeInsightFeed, widget.WidgetType)
			assert.Equal(t, domain.DashboardWidgetSizeMedium, widget.Size)
			assert.Zero(t, widget.MetricDefinitionID)
			assert.Equal(t, 2, widget.OrderIndex)
			assert.Equal(t, usecase.DefaultDashboardConfigJSON, widget.ConfigJSON)
			assert.True(t, widget.IsActive)
			widget.ID = 33
			return widget, nil
		})

	got, err := suite.RecordService.UpsertDashboardWidget(t.Context(), userID, input.UpsertDashboardWidgetCommand{
		ViewID:     viewID,
		WidgetType: domain.DashboardWidgetTypeInsightFeed,
		Size:       domain.DashboardWidgetSizeMedium,
	})

	require.NoError(t, err)
	assert.Equal(t, uint64(33), got.ID)
	assert.Equal(t, domain.DashboardWidgetTypeInsightFeed, got.WidgetType)
	assert.Equal(t, usecase.DefaultDashboardConfigJSON, got.ConfigJSON)
}

func TestService_UpsertDashboardWidget_RequiresMetricDefinitionForNonInsightWidgets(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	_, err := suite.RecordService.UpsertDashboardWidget(t.Context(), 999, input.UpsertDashboardWidgetCommand{
		ViewID:     10,
		WidgetType: domain.DashboardWidgetTypeKPINumber,
		Size:       domain.DashboardWidgetSizeSmall,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, usecase.ErrDashboardMetricDefinitionIDRequired)
}

func TestService_UpsertDashboardWidget_EnforcesLargeWidgetLimit(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	viewID := uint64(10)
	orderIndex := 0

	suite.RecordRepository.EXPECT().
		CountLargeWidgetsInView(gomock.Any(), userID, viewID, (*uint64)(nil)).
		Return(int64(domain.MaxLargeWidgetsPerDashboard), nil)

	_, err := suite.RecordService.UpsertDashboardWidget(t.Context(), userID, input.UpsertDashboardWidgetCommand{
		ViewID:             viewID,
		MetricDefinitionID: 42,
		WidgetType:         domain.DashboardWidgetTypeKPINumber,
		Size:               domain.DashboardWidgetSizeLarge,
		OrderIndex:         &orderIndex,
		ConfigJSON:         `{"layoutVersion":2}`,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, fmt.Sprintf(usecase.ErrDashboardLimitLargeWidgets, domain.MaxLargeWidgetsPerDashboard))
}
