// File: internal/user/core/usecase/get_by_username_test.go
package usecase_test

import (
	"context"
	"testing"

	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
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

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), username).
		Return(u, nil)

	got, err := suite.UserService.GetUserByUsername(context.Background(), username)
	require.NoError(t, err)
	require.Equal(t, u, got)
}

func TestGetUserByUsername_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := "ghost"

	suite.UserRepository.EXPECT().
		GetByUsername(gomock.Any(), username).
		Return(userdomain.User{}, gorm.ErrRecordNotFound)

	got, err := suite.UserService.GetUserByUsername(context.Background(), username)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
}
