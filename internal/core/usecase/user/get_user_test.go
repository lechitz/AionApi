package user_test

import (
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetUserByID_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, userID).
		Return(expectedUser, nil)

	userDomain, err := suite.UserService.GetUserByID(suite.Ctx, userID)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByID_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := setup.DefaultTestUser().ID

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, userID).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserService.GetUserByID(suite.Ctx, userID)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByEmail_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.DefaultTestUser().Email
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, userEmail).
		Return(expectedUser, nil)

	userDomain, err := suite.UserService.GetUserByEmail(suite.Ctx, userEmail)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByEmail_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.DefaultTestUser().Email

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, userEmail).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserService.GetUserByEmail(suite.Ctx, userEmail)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByUsername_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.DefaultTestUser().Username
	expectedUser := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, username).
		Return(expectedUser, nil)

	userDomain, err := suite.UserService.GetUserByUsername(suite.Ctx, username)

	require.NoError(t, err)
	require.Equal(t, expectedUser, userDomain)
}

func TestGetUserByUsername_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.DefaultTestUser().Username

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, username).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserService.GetUserByUsername(suite.Ctx, username)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetAllUsers_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedUsers := []domain.UserDomain{
		setup.DefaultTestUser(),
		setup.DefaultTestUser(),
		setup.DefaultTestUser(),
	}

	suite.UserRepository.EXPECT().
		GetAllUsers(suite.Ctx).
		Return(expectedUsers, nil)

	users, err := suite.UserService.GetAllUsers(suite.Ctx)

	require.NoError(t, err)
	require.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		GetAllUsers(suite.Ctx).
		Return(nil, gorm.ErrRecordNotFound)

	users, err := suite.UserService.GetAllUsers(suite.Ctx)

	require.Error(t, err)
	require.Nil(t, users)
}
