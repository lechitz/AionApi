package auth_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
)

func TestLogin_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := entity.UserDomain{Username: "lechitz"}
	mockUser := entity.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword("hashed", "test123").
		Return(nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, gomock.AssignableToTypeOf(entity.TokenDomain{UserID: 1})).
		Return("token-string", nil)

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "test123")

	require.NoError(t, err)
	require.Equal(t, mockUser, userOut)
	require.Equal(t, "token-string", tokenOut)
}

func TestLogin_UserNotFound(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := entity.UserDomain{Username: "invalid_user"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "invalid_user").
		Return(entity.UserDomain{}, errors.New("not found"))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "123456")

	require.Error(t, err)
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}

func TestLogin_WrongPassword(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := entity.UserDomain{Username: "lechitz"}
	mockUser := entity.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(mockUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword("hashed", "wrongpass").
		Return(errors.New(constants.ErrorToCompareHashAndPassword))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, inputUser, "wrongpass")

	require.Error(t, err)
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
	require.Equal(t, constants.ErrorToCompareHashAndPassword, err.Error())
}

func TestLogin_TokenCreationFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	inputUser := entity.UserDomain{Username: "lechitz"}
	mockUser := entity.UserDomain{ID: 1, Username: "lechitz", Password: "hashed"}

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

	require.Error(t, err)
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
	require.Equal(t, constants.ErrorToCreateToken, err.Error())
}
