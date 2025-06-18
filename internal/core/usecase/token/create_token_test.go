package token_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/lechitz/AionApi/tests/testdata"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// TestCreateToken_NoExistingToken_Success verifies that a token is created and
// saved when none exists for the given user.
func TestCreateToken_NoExistingToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := domain.TokenDomain{UserID: testdata.TestPerfectToken.UserID}

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, input).
		Return("", errors.New("not found"))

	suite.TokenStore.EXPECT().
		Save(suite.Ctx, gomock.AssignableToTypeOf(domain.TokenDomain{UserID: input.UserID})).
		Return(nil)

	token, err := suite.TokenService.CreateToken(suite.Ctx, input)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

// TestCreateToken_WithExistingToken_Success ensures that any existing token is
// removed before creating a new one.
func TestCreateToken_WithExistingToken_Success(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := domain.TokenDomain{UserID: testdata.TestPerfectToken.UserID}

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, input).
		Return("oldtoken", nil)

	suite.TokenStore.EXPECT().
		Delete(suite.Ctx, input).
		Return(nil)

	suite.TokenStore.EXPECT().
		Save(suite.Ctx, gomock.AssignableToTypeOf(domain.TokenDomain{UserID: input.UserID})).
		Return(nil)

	token, err := suite.TokenService.CreateToken(suite.Ctx, input)

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

// TestCreateToken_DeleteFails checks that an error is returned when removing an
// existing token fails.
func TestCreateToken_DeleteFails(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := domain.TokenDomain{UserID: testdata.TestPerfectToken.UserID}

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, input).
		Return("oldtoken", nil)

	suite.TokenStore.EXPECT().
		Delete(suite.Ctx, input).
		Return(errors.New("delete err"))

	token, err := suite.TokenService.CreateToken(suite.Ctx, input)

	require.Error(t, err)
	require.Empty(t, token)
}

// TestCreateToken_SaveFails checks that an error is returned when saving the new
// token fails.
func TestCreateToken_SaveFails(t *testing.T) {
	suite := setup.TokenServiceTest(t, constants.SecretKey)
	defer suite.Ctrl.Finish()

	input := domain.TokenDomain{UserID: testdata.TestPerfectToken.UserID}

	suite.TokenStore.EXPECT().
		Get(suite.Ctx, input).
		Return("", errors.New("not found"))

	suite.TokenStore.EXPECT().
		Save(suite.Ctx, gomock.AssignableToTypeOf(domain.TokenDomain{UserID: input.UserID})).
		Return(errors.New("save err"))

	token, err := suite.TokenService.CreateToken(suite.Ctx, input)

	require.Error(t, err)
	require.Empty(t, token)
}
