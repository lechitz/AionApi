package usecase_test

import (
	"errors"
	"testing"

	authconst "github.com/lechitz/AionApi/internal/feature/auth/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestLogout_Success verifies that Delete(userID) is called and succeeds.
func TestLogout_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), uint64(1)).
		Return(nil)

	err := suite.AuthService.Logout(suite.Ctx, 1)
	require.NoError(t, err)
}

// TestLogout_DeleteTokenFails verifies that a delete failure bubbles up.
func TestLogout_DeleteTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	expected := errors.New("delete error")

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), uint64(1)).
		Return(expected)

	err := suite.AuthService.Logout(suite.Ctx, 1)
	require.Error(t, err)
	require.ErrorContains(t, err, authconst.ErrorToRevokeToken)
}
