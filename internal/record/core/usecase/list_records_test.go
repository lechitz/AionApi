package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/internal/record/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_ListLatest(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		expectedLim int
		repoErr     error
		wantErr     bool
	}{
		{
			name:        "default limit when zero",
			limit:       0,
			expectedLim: 10,
		},
		{
			name:        "uses provided limit",
			limit:       5,
			expectedLim: 5,
		},
		{
			name:        "repository error",
			limit:       3,
			expectedLim: 3,
			repoErr:     errors.New("db error"),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			userID := uint64(1)
			records := []domain.Record{{ID: 1, UserID: userID}}
			if tt.repoErr != nil {
				records = nil
			}

			suite.RecordRepository.EXPECT().
				ListLatest(gomock.Any(), userID, tt.expectedLim).
				Return(records, tt.repoErr)

			result, err := suite.RecordService.ListLatest(suite.Ctx, userID, tt.limit)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), usecase.FailedToListRecords)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_ListByDay(t *testing.T) {
	date := time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)
	userID := uint64(1)
	records := []domain.Record{{ID: 1, UserID: userID}}

	tests := []struct {
		name       string
		cacheErr   error
		cacheHit   bool
		repoErr    error
		wantErr    bool
		expectRepo bool
	}{
		{
			name:     "cache hit",
			cacheHit: true,
		},
		{
			name:       "cache miss repo success",
			cacheErr:   errors.New("cache miss"),
			expectRepo: true,
		},
		{
			name:       "repository error",
			cacheErr:   errors.New("cache miss"),
			repoErr:    errors.New("db error"),
			wantErr:    true,
			expectRepo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			if tt.cacheHit {
				suite.RecordCache.EXPECT().
					GetRecordsByDay(gomock.Any(), userID, date).
					Return(records, nil)
			} else {
				suite.RecordCache.EXPECT().
					GetRecordsByDay(gomock.Any(), userID, date).
					Return(nil, tt.cacheErr)
			}

			if tt.expectRepo {
				repoRecords := records
				if tt.repoErr != nil {
					repoRecords = nil
				}
				suite.RecordRepository.EXPECT().
					ListByDay(gomock.Any(), userID, date).
					Return(repoRecords, tt.repoErr)

				if tt.repoErr == nil {
					suite.RecordCache.EXPECT().
						SaveRecordsByDay(gomock.Any(), userID, date, records, time.Duration(0)).
						Return(nil).
						AnyTimes()
				}
			}

			result, err := suite.RecordService.ListByDay(suite.Ctx, userID, date)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), usecase.FailedToListRecords)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_ListAllBetween(t *testing.T) {
	userID := uint64(1)
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	records := []domain.Record{{ID: 1, UserID: userID}}

	tests := []struct {
		name        string
		start       time.Time
		end         time.Time
		limit       int
		expectedLim int
		repoErr     error
		wantErr     bool
		wantMsg     string
	}{
		{
			name:    "invalid date range",
			start:   end,
			end:     start,
			limit:   10,
			wantErr: true,
			wantMsg: usecase.StartDateMustBeBeforeEndDate,
		},
		{
			name:        "repository error",
			start:       start,
			end:         end,
			limit:       10,
			expectedLim: 10,
			repoErr:     errors.New("db error"),
			wantErr:     true,
			wantMsg:     usecase.FailedToListRecords,
		},
		{
			name:        "default limit when out of range",
			start:       start,
			end:         end,
			limit:       200,
			expectedLim: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			if !tt.start.After(tt.end) {
				suite.RecordRepository.EXPECT().
					ListAllBetween(gomock.Any(), userID, tt.start, tt.end, tt.expectedLim).
					Return(records, tt.repoErr)
			}

			result, err := suite.RecordService.ListAllBetween(suite.Ctx, userID, tt.start, tt.end, tt.limit)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantMsg)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_ListByCategory_LimitBounds(t *testing.T) {
	userID := uint64(1)
	categoryID := uint64(2)
	records := []domain.Record{{ID: 1, UserID: userID}}

	tests := []struct {
		name        string
		limit       int
		expectedLim int
	}{
		{
			name:        "default limit when zero",
			limit:       0,
			expectedLim: 50,
		},
		{
			name:        "default limit when too large",
			limit:       101,
			expectedLim: 50,
		},
		{
			name:        "uses provided limit",
			limit:       20,
			expectedLim: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			suite.RecordRepository.EXPECT().
				ListByCategory(gomock.Any(), categoryID, userID, tt.expectedLim, nil, nil).
				Return(records, nil)

			result, err := suite.RecordService.ListByCategory(suite.Ctx, categoryID, userID, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_ListByTag_LimitBounds(t *testing.T) {
	userID := uint64(1)
	tagID := uint64(3)
	records := []domain.Record{{ID: 1, UserID: userID}}

	tests := []struct {
		name        string
		limit       int
		expectedLim int
	}{
		{
			name:        "default limit when zero",
			limit:       0,
			expectedLim: 50,
		},
		{
			name:        "default limit when too large",
			limit:       150,
			expectedLim: 50,
		},
		{
			name:        "uses provided limit",
			limit:       15,
			expectedLim: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			suite.RecordRepository.EXPECT().
				ListByTag(gomock.Any(), tagID, userID, tt.expectedLim).
				Return(records, nil)

			result, err := suite.RecordService.ListByTag(suite.Ctx, tagID, userID, tt.limit)
			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_ListAllUntil_DefaultLimit(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	until := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
	records := []domain.Record{{ID: 1, UserID: userID}}

	suite.RecordRepository.EXPECT().
		ListAllUntil(gomock.Any(), userID, until, 50).
		Return(records, nil)

	result, err := suite.RecordService.ListAllUntil(suite.Ctx, userID, until, 0)
	require.NoError(t, err)
	assert.Equal(t, records, result)
}

func TestService_ListByUser_RepositoryError(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	suite.RecordRepository.EXPECT().
		ListByUser(gomock.Any(), userID, 10, nil, nil).
		Return(nil, errors.New("db error"))

	result, err := suite.RecordService.ListByUser(suite.Ctx, userID, 10, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), usecase.FailedToListRecords)
	assert.Nil(t, result)
}
