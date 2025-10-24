// Package usecase_test contains tests for the auth use cases.ValidateValidate
package usecase_test

import (
	"errors"
	"strconv"
	"testing"
	"time"

	authDomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	userDomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestLogin_Success verifies that a valid username/password yields a
// generated token saved in the cache and returns the expected user + tokens.
func TestLogin_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "test123").
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("token-string", nil)

	// expect Verify for the generated access token and a positive exp
	exp := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().
		Verify("token-string").
		Return(map[string]any{claimskeys.Exp: exp}, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "token-string", Type: commonkeys.TokenTypeAccess}, gomock.Any()).
		Return(nil)

	// The service generates a refresh token as well; expect the provider call but when it's empty we skip saving it
	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(uint64(1)).
		Return("", nil)

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "test123")

	require.NoError(t, err)
	require.Equal(t, mockUser, userOut)
	require.Equal(t, "token-string", accessTokenOut)
	require.Empty(t, refreshTokenOut)
}

// TestLogin_UserNotFound_ReturnsGetUserError ensures we surface the repository error
// when the username lookup fails.
func TestLogin_UserNotFound_ReturnsGetUserError(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "invalid_user").
		Return(userDomain.User{}, errors.New("not found"))

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "invalid_user", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToGetUserByUserName, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, accessTokenOut)
	require.Empty(t, refreshTokenOut)
}

// TestLogin_WrongPassword_ReturnsInvalidCredentials verifies an invalid password
// returns InvalidCredentials and does not generate/save a token.
func TestLogin_WrongPassword_ReturnsInvalidCredentials(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "wrongpass").
		Return(errors.New(usecase.ErrorToCompareHashAndPassword))

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "wrongpass")

	require.Error(t, err)
	require.Equal(t, usecase.InvalidCredentials, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, accessTokenOut)
	require.Empty(t, refreshTokenOut)
}

// TestLogin_ProviderGenerateFails ensures failure on token generation is surfaced.
func TestLogin_ProviderGenerateFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "123456").
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("", errors.New(usecase.ErrorToCreateToken))

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, accessTokenOut)
	require.Empty(t, refreshTokenOut)
}

// TestLogin_SaveTokenFails ensures failure on saving the token is surfaced.
func TestLogin_SaveTokenFails(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "123456").
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("token-string", nil)

	// Expect Verify for the generated access token and a positive exp so Save is attempted
	exp := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify("token-string").Return(map[string]any{claimskeys.Exp: exp}, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "token-string", Type: commonkeys.TokenTypeAccess}, gomock.Any()).
		Return(errors.New(usecase.ErrorToCreateToken))

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, accessTokenOut)
	require.Empty(t, refreshTokenOut)
}

// TestLogin_UserIDZero_ReturnsInvalidCreds verifies that if the user repository
// returns a user with ID=0, we return the same InvalidCredentials error as when
// the password is wrong.
func TestLogin_UserIDZero_ReturnsInvalidCreds(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "ghost").
		Return(userDomain.User{ID: 0, Username: "ghost"}, nil)

	u, tok, rt, err := suite.AuthService.Login(suite.Ctx, "ghost", "irrelevant")

	require.Error(t, err)
	require.Equal(t, usecase.UserNotFoundOrInvalidCredentials, err.Error())
	require.Empty(t, u)
	require.Empty(t, tok)
	require.Empty(t, rt)
}

func TestLogin_Success_WithRefreshToken(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	suite.Hasher.EXPECT().
		Compare("hashed", "test123").
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("access-token", nil)

	// expect Verify for access token
	expA := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify("access-token").Return(map[string]any{claimskeys.Exp: expA}, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "access-token", Type: commonkeys.TokenTypeAccess}, gomock.Any()).
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(uint64(1)).
		Return("refresh-token", nil)

	// expect Verify for refresh token
	expR := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify("refresh-token").Return(map[string]any{claimskeys.Exp: expR}, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "refresh-token", Type: commonkeys.TokenTypeRefresh}, gomock.Any()).
		Return(nil)

	userOut, accessTokenOut, refreshTokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "test123")

	require.NoError(t, err)
	require.Equal(t, mockUser, userOut)
	require.Equal(t, "access-token", accessTokenOut)
	require.Equal(t, "refresh-token", refreshTokenOut)
}

func TestRefreshTokenRenewal_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	refreshToken := "refresh-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"

	// Expect verification of the refresh token to extract userID
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{Key: userID, Token: refreshToken, Type: commonkeys.TokenTypeRefresh}, nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(userID, gomock.AssignableToTypeOf(map[string]any{})).
		Return(newAccessToken, nil)

	// expect Verify for new access token and save
	expNewA := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify(newAccessToken).Return(map[string]any{claimskeys.Exp: expNewA}, nil)
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: userID, Token: newAccessToken, Type: commonkeys.TokenTypeAccess}, gomock.Any()).
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(userID).
		Return(newRefreshToken, nil)

	// expect Verify for new refresh token and save
	expNewR := strconv.FormatInt(time.Now().Add(24*time.Hour).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify(newRefreshToken).Return(map[string]any{claimskeys.Exp: expNewR}, nil)
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: userID, Token: newRefreshToken, Type: commonkeys.TokenTypeRefresh}, gomock.Any()).
		Return(nil)

	accessTokenOut, refreshTokenOut, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.NoError(t, err)
	require.Equal(t, newAccessToken, accessTokenOut)
	require.Equal(t, newRefreshToken, refreshTokenOut)
}
