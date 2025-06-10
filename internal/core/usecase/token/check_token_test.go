package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
)

func TestVerifyToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	require.NoError(t, err)

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return(tokenString, nil)

	parsedUserID, parsedToken, err := suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	require.NoError(t, err)
	require.Equal(t, userID, parsedUserID)
	require.Equal(t, tokenString, parsedToken)
}

func TestCheck_InvalidToken(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	tokenString := "invalidToken"

	_, _, err := suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	require.Error(t, err)
	require.Equal(t, constants.ErrorInvalidToken, err.Error())
}

func TestCheck_TokenMismatch(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	require.NoError(t, err)

	cachedToken := "differentToken"

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return(cachedToken, nil)

	_, _, err = suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	require.Error(t, err)
	require.Equal(t, constants.ErrorTokenMismatch, err.Error())
}

func TestCheck_ErrorToRetrieveTokenFromCache(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	userID := testdata.TestPerfectUser.ID

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	require.NoError(t, err)

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return("", errors.New(constants.ErrorToRetrieveTokenFromCache))

	_, _, err = suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	require.Error(t, err)
	require.Equal(t, constants.ErrorToRetrieveTokenFromCache, err.Error())
}
