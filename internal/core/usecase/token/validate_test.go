// internal/core/usecase/token/validate_test.go
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

func TestValidate_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 123
	const raw = "Bearer abc.def.ghi"

	suite.TokenProvider.EXPECT().
		Verify(gomock.Any(), "abc.def.ghi").
		Return(map[string]any{"userID": "123"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Token{Key: userID, Value: "abc.def.ghi"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userID, uid)
	require.Equal(t, "123", claims["userID"])
}

func TestValidate_InvalidToken(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer bad.token"
	suite.TokenProvider.EXPECT().
		Verify(gomock.Any(), "bad.token").
		Return(nil, errors.New(constants.ErrorInvalidToken))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
}

func TestValidate_Mismatch(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 5
	const raw = "Bearer abc.def.ghi"

	suite.TokenProvider.EXPECT().
		Verify(gomock.Any(), "abc.def.ghi").
		Return(map[string]any{"userID": "5"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Token{Key: userID, Value: "zzz.yyy.xxx"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Contains(t, err.Error(), constants.ErrorTokenMismatch)
}
