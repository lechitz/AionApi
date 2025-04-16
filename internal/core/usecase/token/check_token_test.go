package token_test

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/infrastructure/security"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerifyToken_Success(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	config.Setting.SecretKey = constants.SecretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	assert.NoError(t, err)

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return(tokenString, nil)

	parsedUserID, parsedToken, err := suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
	assert.Equal(t, tokenString, parsedToken)
}

func TestCheck_InvalidToken(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	config.Setting.SecretKey = constants.SecretKey
	tokenString := "invalidToken"

	_, _, err := suite.TokenService.VerifyToken(context.Background(), tokenString)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrorInvalidToken, err.Error())
}

func TestCheck_TokenMismatch(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	secretKey := constants.SecretKey
	config.Setting.SecretKey = secretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	assert.NoError(t, err)

	cachedToken := "differentToken"

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return(cachedToken, nil)

	_, _, err = suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrorTokenMismatch, err.Error())
}

func TestCheck_ErrorToRetrieveTokenFromCache(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	secretKey := constants.SecretKey
	config.Setting.SecretKey = secretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, constants.SecretKey)
	assert.NoError(t, err)

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return("", errors.New(constants.ErrorToRetrieveTokenFromCache))

	_, _, err = suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrorToRetrieveTokenFromCache, err.Error())
}
