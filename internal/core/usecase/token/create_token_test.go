// Package token_test contains unit tests for the token use case.
package token_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// --- SUITE TESTE TOKEN.
func TestCreateToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 123}

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return("", nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TokenDomain{})).
		Return(nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessTokenCreated, commonkeys.UserID, gomock.Any())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestCreateToken_TokenAlreadyExists_DeleteSuccess(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 99}
	existingToken := "exists"

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return(existingToken, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), tokenDomain).
		Return(nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TokenDomain{})).
		Return(nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessTokenCreated, commonkeys.UserID, gomock.Any())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.NoError(t, err)
	require.NotEmpty(t, token) // Só garante que não está vazio (JWT sempre muda)
}

func TestCreateToken_ErrorToGetToken(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 1}
	expectedErr := errors.New(constants.ErrorToGetToken)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return("", expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetToken, commonkeys.Error, expectedErr.Error())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.Error(t, err)
	require.Empty(t, token)
	require.Equal(t, expectedErr, err)
}

func TestCreateToken_ErrorToDeleteToken(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 1}
	existingToken := "existing"
	expectedErr := errors.New(constants.ErrorToDeleteToken)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return(existingToken, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), tokenDomain).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToDeleteToken, commonkeys.Error, expectedErr.Error())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.Error(t, err)
	require.Empty(t, token)
	require.Equal(t, expectedErr, err)
}

func TestCreateToken_ErrorToAssignToken(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 1}
	expectedErr := errors.New(constants.ErrorToAssignToken)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return("", nil)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToAssignToken, commonkeys.Error, expectedErr.Error())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.Error(t, err)
	require.Empty(t, token)
	require.Equal(t, expectedErr, err)
}

func TestCreateToken_ErrorToSaveToken(t *testing.T) {
	suite := setup.TokenServiceTest(t, config.Secret{Key: "secretKey"})
	defer suite.Ctrl.Finish()

	tokenDomain := domain.TokenDomain{UserID: 1}
	expectedErr := errors.New(constants.ErrorToSaveToken)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), tokenDomain).
		Return("", nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), gomock.AssignableToTypeOf(domain.TokenDomain{})).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToSaveToken, commonkeys.Error, expectedErr.Error())

	token, err := suite.TokenService.CreateToken(context.Background(), tokenDomain)
	require.Error(t, err)
	require.Empty(t, token)
	require.Equal(t, expectedErr, err)
}
