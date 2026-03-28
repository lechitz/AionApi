// File: internal/user/core/usecase/soft_delete_user_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/aion-api/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	userdomain "github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/usecase"
	"github.com/lechitz/aion-api/tests/setup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestSoftDeleteUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	user := userdomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(nil)
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(nil)

	suite.UserRepository.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(nil)

	suite.UserCache.EXPECT().
		DeleteUser(gomock.Any(), userID, user.Username, user.Email).
		Return(nil)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expected := errors.New("error to soft delete user")
	user := userdomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(nil)
	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeRefresh).
		Return(nil)
	suite.UserRepository.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)
	require.Error(t, err)
	require.ErrorContains(t, err, expected.Error())
}

func TestSoftDeleteUser_ErrorToDeleteToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	expected := errors.New(sharederrors.ErrMsgDeleteToken)
	user := userdomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)
	assert.EqualError(t, err, expected.Error())
}

func TestSoftDeleteUser_ContextCancelled(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(1)
	ctx, cancel := context.WithCancel(suite.Ctx)
	cancel()
	user := userdomain.User{
		ID:       userID,
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(context.Canceled)

	err := suite.UserService.SoftDeleteUser(ctx, userID)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestSoftDeleteUser_ErrorToDeleteToken_UnknownError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(2)
	expected := errors.New("unexpected token error")
	user := userdomain.User{
		ID:       userID,
		Username: "testuser2",
		Email:    "test2@example.com",
	}

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(user, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, commonkeys.TokenTypeAccess).
		Return(expected)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)
	assert.EqualError(t, err, expected.Error())
}

// TestSoftDeleteUser_SentinelError validates that SoftDeleteUser returns wrapped ErrSoftDeleteUser.
func TestSoftDeleteUser_SentinelError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := uint64(123)
	dbErr := errors.New("deletion failed")

	suite.UserRepository.EXPECT().
		GetByID(gomock.Any(), userID).
		Return(userdomain.User{ID: userID, Username: "test", Email: "test@example.com"}, nil)

	suite.TokenStore.EXPECT().
		Delete(gomock.Any(), userID, gomock.Any()).
		Return(nil).Times(2)

	suite.UserRepository.EXPECT().
		SoftDelete(gomock.Any(), userID).
		Return(dbErr)

	err := suite.UserService.SoftDeleteUser(suite.Ctx, userID)

	require.Error(t, err)
	require.ErrorIs(t, err, usecase.ErrSoftDeleteUser, "error should wrap ErrSoftDeleteUser sentinel error")
	require.ErrorContains(t, err, "deletion failed")
}
