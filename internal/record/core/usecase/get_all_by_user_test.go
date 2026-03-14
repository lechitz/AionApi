package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/record/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

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
