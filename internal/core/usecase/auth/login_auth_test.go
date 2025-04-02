package auth_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestLogin_Success(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	ctx := suite.Ctx
	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ComparePasswords("hashed", "123456").
		Return(nil)

	suite.TokenService.EXPECT().
		Create(ctx, domain.TokenDomain{UserID: 1}).
		Return("token123", nil)

	suite.TokenService.EXPECT().
		Save(ctx, domain.TokenDomain{UserID: 1, Token: "token123"}).
		Return(nil)

	userOut, tokenOut, err := suite.AuthService.Login(ctx, inputUser, "123456")

	assert.NoError(t, err)
	assert.Equal(t, mockUser, userOut)
	assert.Equal(t, "token123", tokenOut)
}

func TestLogin_UserNotFound(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	ctx := suite.Ctx
	inputUser := domain.UserDomain{Username: "invalid_user"}

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, "invalid_user").
		Return(domain.UserDomain{}, errors.New("not found"))

	userOut, tokenOut, err := suite.AuthService.Login(ctx, inputUser, "123456")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
}

func TestLogin_WrongPassword(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	ctx := suite.Ctx
	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ComparePasswords("hashed", "wrongpass").
		Return(errors.New(constants.ErrorToCompareHashAndPassword))

	userOut, tokenOut, err := suite.AuthService.Login(ctx, inputUser, "wrongpass")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
	assert.Equal(t, constants.ErrorToCompareHashAndPassword, err.Error())
}

func TestLogin_TokenCreationFails(t *testing.T) {
	suite := setup.SetupAuthServiceTest(t)
	defer suite.Ctrl.Finish()

	ctx := suite.Ctx
	inputUser := domain.UserDomain{Username: "lechitz"}
	mockUser := domain.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepo.EXPECT().GetUserByUsername(ctx, "lechitz").Return(mockUser, nil)
	suite.PasswordHasher.EXPECT().ComparePasswords("hashed", "123456").Return(nil)
	suite.TokenService.EXPECT().Create(ctx, domain.TokenDomain{UserID: 1}).Return("", errors.New(constants.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(ctx, inputUser, "123456")

	assert.Error(t, err)
	assert.Empty(t, userOut)
	assert.Empty(t, tokenOut)
	assert.Equal(t, constants.ErrorToCreateToken, err.Error())
}
