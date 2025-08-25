// internal/core/usecase/token/revoke_test.go
package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRevoke_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 11

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(nil)

	err := suite.TokenService.Revoke(suite.Ctx, userID)
	require.NoError(t, err)
}

func TestRevoke_DeleteError(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 11
	exp := errors.New(constants.ErrorToDeleteToken)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(exp)

	err := suite.TokenService.Revoke(suite.Ctx, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), constants.ErrorToDeleteToken)
}
