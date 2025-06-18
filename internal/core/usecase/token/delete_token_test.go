package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

// TestDeleteToken_Success verifies that deleting a token succeeds.
func TestDeleteToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := testdata.TestPerfectToken

	suite.TokenStore.EXPECT().
		Delete(suite.Ctx, input).
		Return(nil)

	err := suite.TokenService.Delete(suite.Ctx, input)

	require.NoError(t, err)
}

// TestDeleteToken_Error ensures repository delete errors are returned.
func TestDeleteToken_Error(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := testdata.TestPerfectToken

	suite.TokenStore.EXPECT().
		Delete(suite.Ctx, input).
		Return(errors.New("delete failed"))

	err := suite.TokenService.Delete(suite.Ctx, input)

	require.Error(t, err)
}
