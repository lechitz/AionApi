package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestLogout_Success verifies that Delete(userID) is called and succeeds.
func TestLogout_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	// Expect both access and refresh token deletes
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), uint64(1), commonkeys.TokenTypeAccess).
		Return(nil)
	// refresh
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), uint64(1), commonkeys.TokenTypeRefresh).
		Return(nil)

	err := suite.AuthService.Logout(suite.Ctx, 1)
	require.NoError(t, err)
}

// TestLogout_DeleteTokenFails verifies that a delete failure bubbles up.
func TestLogout_DeleteTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := uint64(7)

	// Simulate failure on access token delete
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), uid, commonkeys.TokenTypeAccess).
		Return(errors.New("delete error"))

	err := suite.AuthService.Logout(suite.Ctx, uid)

	// Align with current message emitted by use case ("error to delete token")
	require.ErrorContains(t, err, "error to delete token")
}
