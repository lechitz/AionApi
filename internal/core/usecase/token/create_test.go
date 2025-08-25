// internal/core/usecase/token/create_test.go
package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 123
	const jwt = "jwt-token-xyz"

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), userID).
		Return(jwt, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), domain.Token{Key: userID, Value: jwt}).
		Return(nil)

	out, err := suite.TokenService.Create(suite.Ctx, userID)
	require.NoError(t, err)
	require.Equal(t, userID, out.Key)
	require.Equal(t, jwt, out.Value)
}

func TestCreate_ProviderError(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 42
	exp := errors.New(constants.ErrorToAssignToken)

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), userID).
		Return("", exp)

	out, err := suite.TokenService.Create(suite.Ctx, userID)
	require.Error(t, err)
	require.Empty(t, out.Value)
	require.Equal(t, uint64(0), out.Key)
}

func TestCreate_EmptyTokenGenerated(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 99

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), userID).
		Return("", nil)

	out, err := suite.TokenService.Create(suite.Ctx, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), constants.ErrEmptyTokenGenerated)
	require.Equal(t, uint64(0), out.Key)
	require.Empty(t, out.Value)
}

func TestCreate_SaveError(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 7
	const jwt = "jwt"
	exp := errors.New(constants.ErrorToSaveToken)

	suite.TokenProvider.EXPECT().
		Generate(gomock.Any(), userID).
		Return(jwt, nil)

	suite.TokenStore.EXPECT().
		Save(gomock.Any(), domain.Token{Key: userID, Value: jwt}).
		Return(exp)

	out, err := suite.TokenService.Create(suite.Ctx, userID)
	require.Error(t, err)
	require.Contains(t, err.Error(), constants.ErrorToSaveToken)
	require.Equal(t, uint64(0), out.Key)
	require.Empty(t, out.Value)
}
