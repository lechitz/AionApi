// Package user_test contains tests for user use cases.
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSoftDeleteUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(gomock.Any(), userID).
		Return(nil)
	suite.TokenService.EXPECT().
		Delete(gomock.Any(), domain.TokenDomain{UserID: userID}).
		Return(nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUserSoftDeleted, commonkeys.UserID, gomock.Any())

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(gomock.Any(), userID).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToSoftDeleteUser, commonkeys.Error, gomock.Any())

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSoftDeleteUser_ErrorToDeleteToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expectedErr := errors.New(sharederrors.ErrMsgDeleteToken)

	suite.UserRepository.EXPECT().
		SoftDeleteUser(gomock.Any(), userID).
		Return(nil)
	suite.TokenService.EXPECT().
		Delete(gomock.Any(), domain.TokenDomain{UserID: userID}).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), sharederrors.ErrMsgDeleteToken, commonkeys.Error, gomock.Any())

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSoftDeleteUser_ContextCancelled(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	expectedErr := context.Canceled

	suite.UserRepository.EXPECT().
		SoftDeleteUser(gomock.Any(), userID).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToSoftDeleteUser, commonkeys.Error, expectedErr.Error())

	err := suite.UserService.SoftDeleteUser(ctx, userID)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSoftDeleteUser_ErrorToDeleteToken_UnknownError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(2)
	expectedErr := errors.New("unexpected token error")

	suite.UserRepository.EXPECT().
		SoftDeleteUser(gomock.Any(), userID).
		Return(nil)

	suite.TokenService.EXPECT().
		Delete(gomock.Any(), domain.TokenDomain{UserID: userID}).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), sharederrors.ErrMsgDeleteToken, commonkeys.Error, expectedErr.Error())

	err := suite.UserService.SoftDeleteUser(context.Background(), userID)
	assert.EqualError(t, err, expectedErr.Error())
}
