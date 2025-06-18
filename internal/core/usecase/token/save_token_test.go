package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

// TestSaveToken_Success validates that a token is stored successfully.
func TestSaveToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := testdata.TestPerfectToken

	suite.TokenStore.EXPECT().
		Save(suite.Ctx, input).
		Return(nil)

	err := suite.TokenService.Save(suite.Ctx, input)

	require.NoError(t, err)
}

// TestSaveToken_Error ensures that repository errors are propagated.
func TestSaveToken_Error(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := testdata.TestPerfectToken

	suite.TokenStore.EXPECT().
		Save(suite.Ctx, input).
		Return(errors.New("save failed"))

	err := suite.TokenService.Save(suite.Ctx, input)

	require.Error(t, err)
}
