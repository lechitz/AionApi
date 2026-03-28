// File: internal/user/core/usecase/get_user_stats_test.go
package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetUserStats_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(42)
	expectedStats := domain.UserStats{
		TotalRecords:     150,
		TotalCategories:  10,
		TotalTags:        25,
		RecordsThisWeek:  15,
		RecordsThisMonth: 48,
		MostUsedCategory: &domain.UsageCount{
			ID:    1,
			Name:  "Work",
			Count: 45,
		},
		MostUsedTag: &domain.UsageCount{
			ID:    3,
			Name:  "Important",
			Count: 30,
		},
	}

	suite.UserRepository.EXPECT().
		GetUserStats(gomock.Any(), userID).
		Return(expectedStats, nil)

	stats, err := suite.UserService.GetUserStats(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, expectedStats.TotalRecords, stats.TotalRecords)
	require.Equal(t, expectedStats.TotalCategories, stats.TotalCategories)
	require.Equal(t, expectedStats.TotalTags, stats.TotalTags)
	require.Equal(t, expectedStats.RecordsThisWeek, stats.RecordsThisWeek)
	require.Equal(t, expectedStats.RecordsThisMonth, stats.RecordsThisMonth)
	require.NotNil(t, stats.MostUsedCategory)
	require.Equal(t, "Work", stats.MostUsedCategory.Name)
	require.NotNil(t, stats.MostUsedTag)
	require.Equal(t, "Important", stats.MostUsedTag.Name)
}

func TestGetUserStats_EmptyStats(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(99)
	emptyStats := domain.UserStats{
		TotalRecords:     0,
		TotalCategories:  0,
		TotalTags:        0,
		RecordsThisWeek:  0,
		RecordsThisMonth: 0,
		MostUsedCategory: nil,
		MostUsedTag:      nil,
	}

	suite.UserRepository.EXPECT().
		GetUserStats(gomock.Any(), userID).
		Return(emptyStats, nil)

	stats, err := suite.UserService.GetUserStats(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, 0, stats.TotalRecords)
	require.Equal(t, 0, stats.TotalCategories)
	require.Equal(t, 0, stats.TotalTags)
	require.Nil(t, stats.MostUsedCategory)
	require.Nil(t, stats.MostUsedTag)
}

func TestGetUserStats_RepositoryError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	dbErr := errors.New("database connection failed")

	suite.UserRepository.EXPECT().
		GetUserStats(gomock.Any(), userID).
		Return(domain.UserStats{}, dbErr)

	stats, err := suite.UserService.GetUserStats(suite.Ctx, userID)

	require.Error(t, err)
	require.ErrorContains(t, err, "database connection failed")
	require.Equal(t, domain.UserStats{}, stats)
}

func TestGetUserStats_PartialData(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(7)
	partialStats := domain.UserStats{
		TotalRecords:     5,
		TotalCategories:  2,
		TotalTags:        0,
		RecordsThisWeek:  3,
		RecordsThisMonth: 5,
		MostUsedCategory: &domain.UsageCount{
			ID:    2,
			Name:  "Personal",
			Count: 3,
		},
		MostUsedTag: nil, // No tags used yet
	}

	suite.UserRepository.EXPECT().
		GetUserStats(gomock.Any(), userID).
		Return(partialStats, nil)

	stats, err := suite.UserService.GetUserStats(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, 5, stats.TotalRecords)
	require.NotNil(t, stats.MostUsedCategory)
	require.Equal(t, "Personal", stats.MostUsedCategory.Name)
	require.Nil(t, stats.MostUsedTag)
}
