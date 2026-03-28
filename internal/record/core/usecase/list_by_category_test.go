package usecase_test

import (
	"testing"

	"github.com/lechitz/aion-api/internal/record/core/domain"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
