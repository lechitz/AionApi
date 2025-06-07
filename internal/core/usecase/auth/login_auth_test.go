package auth_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword("hashed", "test123").
		Return(nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, gomock.AssignableToTypeOf(domain.TokenDomain{UserID: 1})).
		Return("token-string", nil)

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "test123")

	assert.NoError(t, err)
	assert.Equal(t, mockUser, userOut)
	assert.Equal(t, "token-string", tokenOut)
}

func TestLogin_UserNotFound(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := domain.UserDomain{Username: "invalid_user"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "invalid_user").
		Return(domain.UserDomain{}, errors.New("not found"))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "123456")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
}

func TestLogin_WrongPassword(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword("hashed", "wrongpass").
		Return(errors.New(constants.ErrorToCompareHashAndPassword))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
	assert.Equal(t, constants.ErrorToCompareHashAndPassword, err.Error())
}

func TestLogin_TokenCreationFails(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword("hashed", "123456").
		Return(nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, gomock.Any()).
		Return("", errors.New(constants.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "123456")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
	assert.Equal(t, constants.ErrorToCreateToken, err.Error())
}
