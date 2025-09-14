// Package usecase_test contains tests for the auth use cases.
package usecase_test

import (
	"errors"
	"testing"

	authDomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	userDomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestLogin_Success verifies that a valid username/password yields a
// generated token saved in the cache and returns the expected user + token.
func TestLogin_Success(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	mockUser := userDomain.User{ID: 1, Username: "lechitz", Password: "hashed"}

	// 1) lookup user
	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "lechitz").
		Return(mockUser, nil)

	// 2) compare password (hashed vs plain)
	suite.Hasher.EXPECT().
		Compare("hashed", "test123").
		Return(nil)

	// 3) generate token WITH claims (new contract)
	suite.TokenProvider.EXPECT().
		GenerateWithClaims(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("token-string", nil)

	// 4) persist token
	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "token-string"}).
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
		Return(userDomain.User{}, errors.New("not found"))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "invalid_user", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToGetUserByUserName, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
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

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "wrongpass")

	require.Error(t, err)
	require.Equal(t, usecase.InvalidCredentials, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
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
		GenerateWithClaims(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("", errors.New(usecase.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
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
		GenerateWithClaims(uint64(1), gomock.AssignableToTypeOf(map[string]any{})).
		Return("token-string", nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), authDomain.Auth{Key: 1, Token: "token-string"}).
		Return(errors.New(usecase.ErrorToCreateToken))

	userOut, tokenOut, err := suite.AuthService.Login(suite.Ctx, "lechitz", "123456")

	require.Error(t, err)
	require.Equal(t, usecase.ErrorToCreateToken, err.Error())
	require.Empty(t, userOut)
	require.Empty(t, tokenOut)
}

// TestLogin_UserIDZero_ReturnsInvalidCreds cobre o ramo em que o repositório
// retorna sucesso (err == nil), mas o usuário tem ID == 0. Deve falhar com
// UserNotFoundOrInvalidCredentials, sem chamar hasher/token.
func TestLogin_UserIDZero_ReturnsInvalidCreds(t *testing.T) {
	suite := setup.AuthServiceTest(t)
	defer suite.Ctrl.Finish()

	// Repositório "encontrou", mas o ID é zero => usuário inexistente para o caso de uso.
	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), "ghost").
		Return(userDomain.User{ID: 0, Username: "ghost"}, nil)

	// IMPORTANTE: não deve haver chamadas a Hasher/TokenProvider/TokenStore.
	// Se seu setup tiver EXPECT() default nesses mocks, remova-os para este teste
	// ou use .AnyTimes() nos demais testes, nunca aqui.

	u, tok, err := suite.AuthService.Login(suite.Ctx, "ghost", "irrelevant")

	require.Error(t, err)
	require.Equal(t, usecase.UserNotFoundOrInvalidCredentials, err.Error())
	require.Empty(t, u)
	require.Empty(t, tok)
}
