// Package usecase_test contains tests for the auth use cases.
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

func TestRefreshTokenRenewal_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	refreshToken := "refresh-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"

	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{Key: userID, Token: refreshToken, Type: commonkeys.TokenTypeRefresh}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(authDomain.Auth{}, nil)

	// Mock user cache to return user data (needed for username, email, name in JWT claims)
	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(userDomain.User{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
			Name:     "Test User",
		}, nil)

	suite.RolesReader.EXPECT().
		GetRolesByUserID(gomock.Any(), userID).
		Return([]string{"user"}, nil)

	suite.TokenProvider.EXPECT().
		GenerateAccessToken(userID, gomock.AssignableToTypeOf(map[string]any{})).
		Return(newAccessToken, nil)

	expNewA := strconv.FormatInt(time.Now().Add(3*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().Verify(newAccessToken).Return(map[string]any{claimskeys.Exp: expNewA}, nil)
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: userID, Token: newAccessToken, Type: commonkeys.TokenTypeAccess}, gomock.Any()).
		Return(nil)

	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(userID).
		Return(newRefreshToken, nil)

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

func TestRefreshTokenRenewal_InvalidRefreshToken(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	invalidToken := "invalid-token"

	// Mock: Verify returns error
	suite.TokenProvider.EXPECT().
		Verify(invalidToken).
		Return(nil, errors.New("invalid token signature"))

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, invalidToken)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrInvalidRefreshToken)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestRefreshTokenRenewal_TokenNotInStore(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	refreshToken := "valid-but-not-stored-token" //nolint:gosec // Test token, not a credential

	// Mock: Verify token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Token not found in store
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{}, errors.New("not found"))

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrInvalidRefreshToken)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestRefreshTokenRenewal_TokenMismatch(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	refreshToken := "token-in-request"
	storedToken := "different-token-in-store"

	// Mock: Verify token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Stored token is different (mismatch)
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{
			Key:   userID,
			Token: storedToken,
			Type:  commonkeys.TokenTypeRefresh,
		}, nil)

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrInvalidRefreshToken)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestRefreshTokenRenewal_UserNotFound(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(999)
	refreshToken := "valid-refresh-token"

	// Mock: Verify refresh token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Get stored refresh token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{
			Key:   userID,
			Token: refreshToken,
			Type:  commonkeys.TokenTypeRefresh,
		}, nil)

	// Mock: Get old access token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(authDomain.Auth{}, nil)

	// Mock: User not found in cache
	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(userDomain.User{}, errors.New("not found"))

	// Mock: User not found in repository
	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(userDomain.User{}, errors.New("user not found"))

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrTokenCreation)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestRefreshTokenRenewal_FailToGenerateAccessToken(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	refreshToken := "valid-refresh-token"

	// Mock: Verify refresh token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Get stored refresh token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{
			Key:   userID,
			Token: refreshToken,
			Type:  commonkeys.TokenTypeRefresh,
		}, nil)

	// Mock: Get old access token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(authDomain.Auth{}, nil)

	// Mock: Get user
	mockUser := userDomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}
	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(mockUser, nil)

	// Mock: Get roles
	suite.RolesReader.EXPECT().
		GetRolesByUserID(gomock.Any(), userID).
		Return([]string{"user"}, nil)

	// Mock: Fail to generate access token
	suite.TokenProvider.EXPECT().
		GenerateAccessToken(userID, gomock.Any()).
		Return("", errors.New("token generation failed"))

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrTokenCreation)
	require.Empty(t, access)
	require.Empty(t, refresh)
}

func TestRefreshTokenRenewal_GracePeriodCreated(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	refreshToken := "valid-refresh-token"
	oldAccessToken := "old-access-token"
	newAccessToken := "new-access-token"
	newRefreshToken := "new-refresh-token"

	// Mock: Verify refresh token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Get stored refresh token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{
			Key:   userID,
			Token: refreshToken,
			Type:  commonkeys.TokenTypeRefresh,
		}, nil)

	// Mock: Get old access token (should trigger grace period)
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(authDomain.Auth{
			Key:   userID,
			Token: oldAccessToken,
			Type:  commonkeys.TokenTypeAccess,
		}, nil)

	// Mock: Get user
	mockUser := userDomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}
	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(mockUser, nil)

	// Mock: Get roles
	suite.RolesReader.EXPECT().
		GetRolesByUserID(gomock.Any(), userID).
		Return([]string{"user"}, nil)

	// Mock: Generate new access token
	suite.TokenProvider.EXPECT().
		GenerateAccessToken(userID, gomock.Any()).
		Return(newAccessToken, nil)

	// Mock: Verify new access token (for TTL)
	accessExp := strconv.FormatInt(time.Now().Add(15*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().
		Verify(newAccessToken).
		Return(map[string]any{claimskeys.Exp: accessExp}, nil)

	// Mock: Save new access token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// Mock: SaveWithKey for grace period (this is what we're testing)
	suite.TokenStore.EXPECT().
		SaveWithKey(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ any, key string, auth authDomain.Auth, ttl time.Duration) error {
			// Verify grace period key format
			require.Contains(t, key, "auth:grace:")
			require.Contains(t, key, strconv.FormatUint(userID, 10))
			// Verify TTL is grace period duration (60 seconds)
			require.Equal(t, 60*time.Second, ttl)
			// Verify old token is stored
			require.Equal(t, oldAccessToken, auth.Token)
			return nil
		})

	// Mock: Generate new refresh token
	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(userID).
		Return(newRefreshToken, nil)

	// Mock: Verify new refresh token (for TTL)
	refreshExp := strconv.FormatInt(time.Now().Add(7*24*time.Hour).Unix(), 10)
	suite.TokenProvider.EXPECT().
		Verify(newRefreshToken).
		Return(map[string]any{claimskeys.Exp: refreshExp}, nil)

	// Mock: Save new refresh token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.NoError(t, err)
	require.Equal(t, newAccessToken, access)
	require.Equal(t, newRefreshToken, refresh)
}

func TestRefreshTokenRenewal_UserCacheMiss(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	refreshToken := "valid-refresh-token"

	// Mock: Verify refresh token
	suite.TokenProvider.EXPECT().
		Verify(refreshToken).
		Return(map[string]any{claimskeys.UserID: strconv.FormatUint(userID, 10)}, nil)

	// Mock: Get stored refresh token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(authDomain.Auth{
			Key:   userID,
			Token: refreshToken,
			Type:  commonkeys.TokenTypeRefresh,
		}, nil)

	// Mock: Get old access token
	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(authDomain.Auth{}, nil)

	// Mock: Cache miss
	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(userDomain.User{}, errors.New("cache miss"))

	// Mock: Get user from repository
	mockUser := userDomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}
	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(mockUser, nil)

	// Mock: Save user to cache after retrieval
	suite.UserCache.EXPECT().
		SaveUser(gomock.Any(), mockUser, gomock.Any()).
		Return(nil)

	// Mock: Get roles
	suite.RolesReader.EXPECT().
		GetRolesByUserID(gomock.Any(), userID).
		Return([]string{"user"}, nil)

	// Mock: Generate new access token
	suite.TokenProvider.EXPECT().
		GenerateAccessToken(userID, gomock.Any()).
		Return("new-access-token", nil)

	// Mock: Verify new access token
	accessExp := strconv.FormatInt(time.Now().Add(15*time.Minute).Unix(), 10)
	suite.TokenProvider.EXPECT().
		Verify("new-access-token").
		Return(map[string]any{claimskeys.Exp: accessExp}, nil)

	// Mock: Save new access token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	// Mock: Generate new refresh token
	suite.TokenProvider.EXPECT().
		GenerateRefreshToken(userID).
		Return("new-refresh-token", nil)

	// Mock: Verify new refresh token
	refreshExp := strconv.FormatInt(time.Now().Add(7*24*time.Hour).Unix(), 10)
	suite.TokenProvider.EXPECT().
		Verify("new-refresh-token").
		Return(map[string]any{claimskeys.Exp: refreshExp}, nil)

	// Mock: Save new refresh token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	access, refresh, err := suite.AuthService.RefreshTokenRenewal(suite.Ctx, refreshToken)

	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)
}
