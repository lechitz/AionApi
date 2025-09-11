// File: internal/user/core/usecase/get_by_id_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetUserByID_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	u := userdomain.User{
		ID:       1,
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "f@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), u.ID).
		Return(u, nil)

	got, err := suite.UserService.GetByID(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, u, got)
}

func TestGetUserByID_ErrorGeneric(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	id := uint64(99)
	expected := errors.New("some db failure")

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(userdomain.User{}, expected)

	got, err := suite.UserService.GetByID(context.Background(), id)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
	require.ErrorContains(t, err, "some db failure")
}

func TestGetUserByID_NotFound(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	id := uint64(100)

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(userdomain.User{}, gorm.ErrRecordNotFound)

	got, err := suite.UserService.GetByID(context.Background(), id)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
}
