package usecase_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/record/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_ListRecordsByCategory(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	categoryID := uint64(2)
	records := []domain.Record{{ID: 1, UserID: userID}}

	suite.RecordRepository.EXPECT().
		ListByCategory(gomock.Any(), categoryID, userID, 10, nil, nil).
		Return(records, nil)

	result, err := suite.RecordService.ListRecordsByCategory(suite.Ctx, categoryID, userID, 10)
	require.NoError(t, err)
	assert.Equal(t, records, result)
}
