// File: internal/user/core/usecase/soft_delete_user_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSoftDeleteUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	// Order: authStore.Delete -> userRepository.SoftDelete
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(nil)
	suite.UserRepository.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(nil)

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expected := errors.New("error to soft delete user")

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(nil)
	suite.UserRepository.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expected.Error())
}

func TestSoftDeleteUser_ErrorToDeleteToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expected := errors.New(sharederrors.ErrMsgDeleteToken)

	// If token deletion fails, repository must NOT be called.
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expected.Error())
}

func TestSoftDeleteUser_ContextCancelled(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(context.Canceled)

	err := suite.UserService.SoftDeleteUser(ctx, userID)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSoftDeleteUser_ErrorToDeleteToken_UnknownError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(2)
	expected := errors.New("unexpected token error")

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expected.Error())
}
