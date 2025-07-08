package auth_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestLogout_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	token := "valid.token.value" // #nosec G101
	userID := uint64(1)

	suite.TokenService.EXPECT().
		GetToken(suite.Ctx, token).
		Return(userID, token, nil)

	suite.TokenService.EXPECT().
		Delete(suite.Ctx, domain.TokenDomain{UserID: userID, Token: token}).
		Return(nil)

	err := suite.AuthService.Logout(suite.Ctx, token)

	assert.NoError(t, err)
}

func TestLogout_CheckTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	token := "invalid.token.value"

	suite.TokenService.EXPECT().
		GetToken(suite.Ctx, token).
		Return(uint64(0), "", errors.New(constants.ErrorInvalidToken))

	err := suite.AuthService.Logout(suite.Ctx, token)

	assert.ErrorContains(t, err, constants.ErrorInvalidToken)
}

func TestLogout_DeleteTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	token := "valid.token.value" // #nosec G101
	userID := uint64(1)
	tokenDomain := domain.TokenDomain{UserID: userID, Token: token}
	expectedErr := errors.New("delete error")

	suite.TokenService.EXPECT().
		GetToken(suite.Ctx, token).
		Return(userID, token, nil)

	suite.TokenService.EXPECT().
		Delete(suite.Ctx, tokenDomain).
		Return(expectedErr)

	err := suite.AuthService.Logout(suite.Ctx, token)

	assert.ErrorContains(t, err, expectedErr.Error())
}
