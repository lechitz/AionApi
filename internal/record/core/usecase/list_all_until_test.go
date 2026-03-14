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

func TestService_ListAllUntil_RepositoryError(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	until := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)

	suite.RecordRepository.EXPECT().
		ListAllUntil(gomock.Any(), userID, until, 10).
		Return(nil, errors.New("db error"))

	result, err := suite.RecordService.ListAllUntil(suite.Ctx, userID, until, 10)
	require.Error(t, err)
	assert.Contains(t, err.Error(), usecase.FailedToListRecords)
	assert.Nil(t, result)
}
