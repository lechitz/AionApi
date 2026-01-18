// File: internal/user/core/usecase/get_by_id_test.go
package usecase_test

import (
	"errors"
	"testing"

	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/usecase"
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

	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), u.ID).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), u.ID).
		Return(u, nil)

	suite.UserCache.EXPECT().
		SaveUser(gomock.Any(), u, gomock.Any()).
		Return(nil)

	got, err := suite.UserService.GetByID(suite.Ctx, u.ID)
	require.NoError(t, err)
	require.Equal(t, u, got)
}

func TestGetUserByID_ErrorGeneric(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	id := uint64(99)
	expected := errors.New("some db failure")

	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), id).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(userdomain.User{}, expected)

	got, err := suite.UserService.GetByID(suite.Ctx, id)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
	require.ErrorContains(t, err, "some db failure")
}

func TestGetUserByID_NotFound(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	id := uint64(100)

	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), id).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(userdomain.User{}, gorm.ErrRecordNotFound)

	got, err := suite.UserService.GetByID(suite.Ctx, id)
	require.Error(t, err)
	require.Equal(t, userdomain.User{}, got)
}

// TestGetUserByID_SentinelError validates that GetByID returns wrapped ErrGetSelf.
func TestGetUserByID_SentinelError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	id := uint64(999)
	dbErr := errors.New("database connection failed")

	suite.UserCache.EXPECT().
		GetUserByID(gomock.Any(), id).
		Return(userdomain.User{}, errors.New("cache miss"))

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), id).
		Return(userdomain.User{}, dbErr)

	_, err := suite.UserService.GetByID(suite.Ctx, id)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrGetSelf, "error should wrap ErrGetSelf sentinel error")
	require.ErrorContains(t, err, "database connection failed")
}
