package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
