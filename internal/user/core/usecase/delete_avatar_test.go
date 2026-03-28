package usecase_test

import (
	"errors"
	"testing"

	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestRemoveAvatar_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	expected := setup.DefaultTestUser()
	expected.AvatarURL = nil

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	suite.UserCache.EXPECT().
		DeleteUser(gomock.Any(), uid, expected.Username, expected.Email).
		Return(nil)

	got, err := suite.UserService.RemoveAvatar(suite.Ctx, uid)
	require.NoError(t, err)
	require.Equal(t, expected, got)
	require.Nil(t, got.AvatarURL)
}

func TestRemoveAvatar_ErrorToUpdateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	repoErr := errors.New("update failed")

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(userdomain.User{}, repoErr)

	got, err := suite.UserService.RemoveAvatar(suite.Ctx, uid)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
	require.ErrorIs(t, err, usecase.ErrDeleteAvatar)
	require.ErrorContains(t, err, repoErr.Error())
}
