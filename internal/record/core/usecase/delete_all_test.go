package usecase_test

import (
	"testing"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_DeleteAll(t *testing.T) {
	suite := setup.RecordServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	suite.RecordRepository.EXPECT().
		DeleteAllByUser(gomock.Any(), userID).
		Return(nil)

	err := suite.RecordService.DeleteAll(suite.Ctx, userID)
	require.NoError(t, err)
}
