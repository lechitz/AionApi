package user_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"

	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestSoftDeleteUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(suite.Ctx, userID).
		Return(nil)

	suite.TokenService.EXPECT().
		Delete(suite.Ctx, entity.TokenDomain{UserID: userID}).
		Return(nil)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)

	assert.NoError(t, err)
}

func TestSoftDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(suite.Ctx, userID).
		Return(expectedErr)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)

	assert.EqualError(t, err, expectedErr.Error())
}

func TestSoftDeleteUser_ErrorToDeleteToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedErr := errors.New(constants.ErrorToDeleteToken)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(suite.Ctx, userID).
		Return(nil)

	suite.TokenService.EXPECT().
		Delete(suite.Ctx, entity.TokenDomain{UserID: userID}).
		Return(expectedErr)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)

	assert.EqualError(t, err, expectedErr.Error())
}
