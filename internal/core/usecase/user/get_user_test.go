package user_test

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestGetUserByID_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	expectedUser := setup.TestPerfectUser

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, userID).
		Return(expectedUser, nil)

	userDomain, err := suite.UserSvc.GetUserByID(suite.Ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, userDomain)
}

func TestGetUserByID_Error(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := setup.TestPerfectUser.ID

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, userID).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserSvc.GetUserByID(suite.Ctx, userID)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByEmail_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.TestPerfectUser.Email

	expectedUser := setup.TestPerfectUser

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, userEmail).
		Return(expectedUser, nil)

	userDomain, err := suite.UserSvc.GetUserByEmail(suite.Ctx, userEmail)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, userDomain)
}

func TestGetUserByEmail_Error(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	userEmail := setup.TestPerfectUser.Email

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, userEmail).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserSvc.GetUserByEmail(suite.Ctx, userEmail)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetUserByUsername_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.TestPerfectUser.Username

	expectedUser := setup.TestPerfectUser

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, username).
		Return(expectedUser, nil)

	userDomain, err := suite.UserSvc.GetUserByUsername(suite.Ctx, username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, userDomain)
}

func TestGetUserByUsername_Error(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	username := setup.TestPerfectUser.Username

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, username).
		Return(domain.UserDomain{}, gorm.ErrRecordNotFound)

	userDomain, err := suite.UserSvc.GetUserByUsername(suite.Ctx, username)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, userDomain)
}

func TestGetAllUsers_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	expectedUsers := []domain.UserDomain{
		setup.TestPerfectUser,
		setup.TestPerfectUser,
		setup.TestPerfectUser,
	}

	suite.UserRepository.EXPECT().
		GetAllUsers(suite.Ctx).
		Return(expectedUsers, nil)

	users, err := suite.UserSvc.GetAllUsers(suite.Ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestGetAllUsers_Error(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	suite.UserRepository.EXPECT().
		GetAllUsers(suite.Ctx).
		Return(nil, gorm.ErrRecordNotFound)

	users, err := suite.UserSvc.GetAllUsers(suite.Ctx)

	assert.Error(t, err)
	assert.Nil(t, users)
}
