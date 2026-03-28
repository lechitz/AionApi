// File: internal/user/core/usecase/get_by_username_test.gohttputils
package usecase_test

import (
	"errors"
	"testing"

	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetUserByUsername_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := "lechitz"
	u := userdomain.User{
		ID:       42,
		Name:     "Felipe",
		Username: username,
		Email:    "f@example.com",
	}

	suite.UserCache.EXPECT().
		GetUserByUsername(gomock.Any(), username).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), username).
		Return(u, nil)

	suite.UserCache.EXPECT().
		SaveUser(gomock.Any(), u, gomock.Any()).
		Return(nil)

	got, err := suite.UserService.GetUserByUsername(suite.Ctx, username)
	require.NoError(t, err)
	require.Equal(t, u, got)
}

func TestGetUserByUsername_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := "ghost"

	suite.UserCache.EXPECT().
		GetUserByUsername(gomock.Any(), username).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), username).
		Return(userdomain.User{}, gorm.ErrRecordNotFound)

	got, err := suite.UserService.GetUserByUsername(suite.Ctx, username)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
}

// TestGetUserByUsername_SentinelError validates that GetByUsername returns wrapped ErrGetUserByUsername.
func TestGetUserByUsername_SentinelError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := "nonexistent"
	dbErr := errors.New("user not found")

	suite.UserCache.EXPECT().
		GetUserByUsername(gomock.Any(), username).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), username).
		Return(userdomain.User{}, dbErr)

	_, err := suite.UserService.GetUserByUsername(suite.Ctx, username)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrGetUserByUsername, "error should wrap ErrGetUserByUsername sentinel error")
	require.ErrorContains(t, err, "user not found")
}
