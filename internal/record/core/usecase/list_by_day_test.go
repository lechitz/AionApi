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
