// File: internal/user/core/usecase/update_user_test.go
package usecase_test

import (
	"errors"
	"testing"

	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	userinput "github.com/lechitz/AionApi/internal/user/core/ports/input"
	"github.com/lechitz/AionApi/internal/user/core/usecase"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func strptr(s string) *string { return &s }

func TestUpdateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	cmd := userinput.UpdateUserCommand{
		Name:     strptr("Felipe"),
		Username: strptr(setup.DefaultTestUser().Username),
		Email:    strptr(setup.DefaultTestUser().Email),
	}
	expected := setup.DefaultTestUser()
	expected.Name = "Felipe"

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	got, err := suite.UserService.UpdateUser(suite.Ctx, uid, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestUpdateUser_UpdateOnlyUsername(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	cmd := userinput.UpdateUserCommand{
		Username: strptr("new_username"),
	}
	expected := setup.DefaultTestUser()
	expected.Username = "new_username"

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	got, err := suite.UserService.UpdateUser(suite.Ctx, uid, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestUpdateUser_UpdateOnlyEmail(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	cmd := userinput.UpdateUserCommand{
		Email: strptr("new@email.com"),
	}
	expected := setup.DefaultTestUser()
	expected.Email = "new@email.com"

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	got, err := suite.UserService.UpdateUser(suite.Ctx, uid, cmd)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

func TestUpdateUser_NoFieldsToUpdate(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	var cmd userinput.UpdateUserCommand // all fields nil => HasUpdates() == false

	got, err := suite.UserService.UpdateUser(suite.Ctx, uid, cmd)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
	require.Equal(t, usecase.ErrorNoFieldsToUpdate, err.Error())
}

func TestUpdateUser_ErrorToUpdateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	uid := setup.DefaultTestUser().ID
	cmd := userinput.UpdateUserCommand{
		Name:     strptr("Felipe"),
		Username: strptr(setup.DefaultTestUser().Username),
		Email:    strptr(setup.DefaultTestUser().Email),
	}
	repoErr := errors.New(usecase.ErrorToUpdateUser)

	suite.UserRepository.EXPECT().
		Update(gomock.Any(), uid, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(userdomain.User{}, repoErr)

	got, err := suite.UserService.UpdateUser(suite.Ctx, uid, cmd)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
	require.ErrorContains(t, err, usecase.ErrorToUpdateUser)
}
