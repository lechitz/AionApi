package usecase_test

import (
	"errors"
	"testing"

	domain2 "github.com/lechitz/AionApi/internal/core/user/domain"
	"github.com/lechitz/AionApi/internal/feature/auth/core/domain"
	authconst "github.com/lechitz/AionApi/internal/feature/auth/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestLogin_Success verifies that a valid username/password yields a
// generated token saved in the cache and returns the expected user + token.
func TestLogin_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := domain2.User{ID: 1, Username: "lechitz", Password: "hashed"}

	// 1) lookup user
	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	// 2) compare password (hashed vs plain)
	suite.Hasher.EXPECT().
		Compare("hashed", "test123").
		Return(nil)

	// 3) generate token
	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), uint64(1)).
		Return("token-string", nil)

	// 4) persist token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), domain.Auth{Key: 1, Token: "token-string"}).
		Return(nil)

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "test123")

	require.NoError(t, err)
	require.Equal(t, mockUser, userOut)
	require.Equal(t, "token-string", tokenOut)
}

// TestLogin_UserNotFound_ReturnsGetUserError ensures we surface the repository error
// when the username lookup fails.
func TestLogin_UserNotFound_ReturnsGetUserError(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "invalid_user").
		Return(domain2.User{}, errors.New("not found"))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "invalid_user", "123456")

	require.Error(t, err)
	require.Equal(t, authconst.ErrorToGetUserByUserName, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}

// TestLogin_WrongPassword_ReturnsInvalidCredentials verifies an invalid password
// returns InvalidCredentials and does not generate/save a token.
func TestLogin_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := domain2.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "wrongpass").
		Return(errors.New(authconst.ErrorToCompareHashAndPassword))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "wrongpass")

	require.Error(t, err)
	require.Equal(t, authconst.InvalidCredentials, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}

// TestLogin_ProviderGenerateFails ensures failure on token generation is surfaced.
func TestLogin_ProviderGenerateFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := domain2.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "123456").
		Return(nil)

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), uint64(1)).
		Return("", errors.New(authconst.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, authconst.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}

// TestLogin_SaveTokenFails ensures failure on saving the token is surfaced.
func TestLogin_SaveTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := domain2.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "123456").
		Return(nil)

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), uint64(1)).
		Return("token-string", nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), domain.Auth{Key: 1, Token: "token-string"}).
		Return(errors.New(authconst.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, authconst.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}
