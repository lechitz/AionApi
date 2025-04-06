package token_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/infrastructure/security"
	"github.com/lechitz/AionApi/internal/platform/config"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestVerifyToken_Success(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, "secret")
	defer suite.Ctrl.Finish()

	secretKey := "secret"
	config.Setting.SecretKey = secretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, "secret")
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTokenRepo := mockToken.NewMockTokenRepositoryPort(ctrl)
	logger := zaptest.NewLogger(t).Sugar()
	secretKey := "secret"

	tokenSvc := token.NewTokenService(mockTokenRepo, logger, domain.TokenConfig{
		SecretKey: secretKey,
	})

	ctx := domain.ContextControl{}
	tokenString := "invalidToken"

	_, _, err := tokenSvc.VerifyToken(ctx, tokenString)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrorInvalidToken, err.Error())
}

func TestCheck_TokenMismatch(t *testing.T) {
	suite := setup.SetupTokenServiceTest(t, "secret")
	defer suite.Ctrl.Finish()

	secretKey := "secret"
	config.Setting.SecretKey = secretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, "secret")
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
	suite := setup.SetupTokenServiceTest(t, "secret")
	defer suite.Ctrl.Finish()

	secretKey := "secret"
	config.Setting.SecretKey = secretKey

	userID := uint64(1)

	tokenString, err := security.GenerateToken(userID, "secret")
	assert.NoError(t, err)

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, domain.TokenDomain{UserID: userID, Token: tokenString}).
		Return("", errors.New(constants.ErrorToRetrieveTokenFromCache))

	_, _, err = suite.TokenService.VerifyToken(suite.Ctx, tokenString)

	assert.Error(t, err)
	assert.Equal(t, constants.ErrorToRetrieveTokenFromCache, err.Error())
}
