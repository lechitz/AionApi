package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_SearchRecords_Limits(t *testing.T) {
	tests := []struct {
		name          string
		inputLimit    int
		expectedLimit int
	}{
		{
			name:          "default limit",
			inputLimit:    0,
			expectedLimit: 20,
		},
		{
			name:          "cap limit",
			inputLimit:    150,
			expectedLimit: 100,
		},
		{
			name:          "use provided limit",
			inputLimit:    25,
			expectedLimit: 25,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suite := setup.RecordServiceTest(t)
			defer suite.Ctrl.Finish()

			userID := uint64(1)
			filters := domain.SearchFilters{Query: "hello", Limit: tt.inputLimit}
			records := []domain.Record{{ID: 1, UserID: userID}}

			suite.RecordRepository.EXPECT().
				SearchRecords(gomock.Any(), userID, gomock.AssignableToTypeOf(domain.SearchFilters{})).
				DoAndReturn(func(_ context.Context, _ uint64, f domain.SearchFilters) ([]domain.Record, error) {
					assert.Equal(t, tt.expectedLimit, f.Limit)
					return records, nil
				})

			result, err := suite.RecordService.SearchRecords(suite.Ctx, userID, filters)
			require.NoError(t, err)
			assert.Equal(t, records, result)
		})
	}
}

func TestService_SearchRecords_Error(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	filters := domain.SearchFilters{Query: "oops", Limit: 10}

	suite.RecordRepository.EXPECT().
		SearchRecords(gomock.Any(), userID, filters).
		Return(nil, errors.New("search failed"))

	result, err := suite.RecordService.SearchRecords(suite.Ctx, userID, filters)
	require.Error(t, err)
	assert.Nil(t, result)
}
