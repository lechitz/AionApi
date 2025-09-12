package usecase_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestValidate_Success verifies that a valid token yields the expected user ID.
func TestValidate_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 123
	const raw = "Bearer abc.def.ghi"

	suite.AuthProvider.EXPECT().
		Verify("abc.def.ghi").
		Return(map[string]any{"userID": "123"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Auth{Key: userID, Token: "abc.def.ghi"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.NoError(t, err)
	require.Equal(t, userID, uid)
	require.Equal(t, "123", claims["userID"])
}

// TestValidate_InvalidToken verifies that an invalid token is surfaced.
func TestValidate_InvalidToken(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const raw = "Bearer bad.token"
	suite.AuthProvider.EXPECT().
		Verify("bad.token").
		Return(nil, errors.New(usecase.ErrorInvalidToken))

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
}

// TestValidate_Mismatch verifies that a mismatch between the token and the stored token is surfaced.
func TestValidate_Mismatch(t *testing.T) {
	suite := setup.TokenServiceTest(t)
	defer suite.Ctrl.Finish()

	const userID uint64 = 5
	const raw = "Bearer abc.def.ghi"

	suite.AuthProvider.EXPECT().
		Verify("abc.def.ghi").
		Return(map[string]any{"userID": "5"}, nil)

	suite.TokenStore.EXPECT().
		Get(gomock.Any(), userID).
		Return(domain.Auth{Key: userID, Token: "zzz.yyy.xxx"}, nil)

	uid, claims, err := suite.TokenService.Validate(suite.Ctx, raw)
	require.Error(t, err)
	require.Equal(t, uint64(0), uid)
	require.Nil(t, claims)
	require.Contains(t, err.Error(), usecase.ErrorTokenMismatch)
}
