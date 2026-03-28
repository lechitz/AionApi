package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/internal/record/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
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
