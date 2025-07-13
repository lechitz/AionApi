package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestGetUserByID_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(expectedUser, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUserRetrieved, commonkeys.UserID, gomock.Any())

	userDomain, err := suite.UserService.GetUserByID(context.Background(), userID)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByID_ErrorGeneric(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	expectedErr := errors.New("some db failure")

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByID, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByID(context.Background(), userID)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
	require.Contains(t, err.Error(), "some db failure")
}

func TestGetUserByID_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := setup.DefaultTestUser().ID
	expectedErr := gorm.ErrRecordNotFound

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByID, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByID(context.Background(), userID)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByEmail_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.DefaultTestUser().Email
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), userEmail).
		Return(expectedUser, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUserRetrieved, commonkeys.UserID, gomock.Any())

	userDomain, err := suite.UserService.GetUserByEmail(context.Background(), userEmail)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByEmail_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.DefaultTestUser().Email
	expectedErr := gorm.ErrRecordNotFound

	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), userEmail).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByEmail, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByEmail(context.Background(), userEmail)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByUsername_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.DefaultTestUser().Username
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), username).
		Return(expectedUser, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUserRetrieved, commonkeys.UserID, gomock.Any())

	userDomain, err := suite.UserService.GetUserByUsername(context.Background(), username)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByUsername_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.DefaultTestUser().Username
	expectedErr := gorm.ErrRecordNotFound

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), username).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByUsername, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByUsername(context.Background(), username)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetAllUsers_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedUsers := []domain.UserDomain{
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

	var expectedUsers []domain.UserDomain

	suite.UserRepository.EXPECT().
		GetAllUsers(gomock.Any()).
		Return(expectedUsers, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUsersRetrieved, commonkeys.Users, "0")

	users, err := suite.UserService.GetAllUsers(context.Background())

	require.NoError(t, err)
	require.Empty(t, users)
}
