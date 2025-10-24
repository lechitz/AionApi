// File: internal/user/core/usecase/list_all_test.go
package usecase_test

import (
	"errors"
	"testing"

	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListAll_Success_WithUsers(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	users := []userdomain.User{
		{ID: 1, Name: "U1", Username: "u1", Email: "u1@ex.com"},
		{ID: 2, Name: "U2", Username: "u2", Email: "u2@ex.com"},
	}

	suite.UserRepository.EXPECT().
		ListAll(gomock.Any()).
		Return(users, nil)

	got, err := suite.UserService.ListAll(suite.Ctx)
	require.NoError(t, err)
	require.Equal(t, users, got)
}

func TestListAll_Success_EmptySlice(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		ListAll(gomock.Any()).
		Return([]userdomain.User{}, nil)

	got, err := suite.UserService.ListAll(suite.Ctx)
	require.NoError(t, err)
	require.Empty(t, got)
}

func TestListAll_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	expected := errors.New("db failure")

	suite.UserRepository.EXPECT().
		ListAll(gomock.Any()).
		Return(nil, expected)

	got, err := suite.UserService.ListAll(suite.Ctx)
	require.Error(t, err)
	require.Empty(t, got) // accept []domain.User{} or nil
	require.ErrorContains(t, err, "db failure")
}
