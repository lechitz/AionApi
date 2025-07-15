// Package user_test contains tests for user use cases.
package user_test

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetAllUsers_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedUsers := []domain.User{
		setup.DefaultTestUser(),
		setup.DefaultTestUser(),
	}

	suite.UserRepository.EXPECT().
		GetAllUsers(gomock.Any()).
		Return(expectedUsers, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUsersRetrieved, commonkeys.Users, gomock.Any())

	users, err := suite.UserService.GetAllUsers(context.Background())

	require.NoError(t, err)
	require.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedErr := gorm.ErrRecordNotFound

	suite.UserRepository.EXPECT().
		GetAllUsers(gomock.Any()).
		Return(nil, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetAllUsers, commonkeys.Error, expectedErr.Error())

	users, err := suite.UserService.GetAllUsers(context.Background())

	require.Error(t, err)
	require.Nil(t, users)
}

func TestGetAllUsers_EmptyResult(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	var expectedUsers []domain.User

	suite.UserRepository.EXPECT().
		GetAllUsers(gomock.Any()).
		Return(expectedUsers, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUsersRetrieved, commonkeys.Users, "0")

	users, err := suite.UserService.GetAllUsers(context.Background())

	require.NoError(t, err)
	require.Empty(t, users)
}
