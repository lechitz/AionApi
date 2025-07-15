// Package user_test contains tests for user use cases.
package user_test

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
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
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByID, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByID(context.Background(), userID)

	require.Error(t, err)
	require.Equal(t, domain.User{}, userDomain)
	require.Contains(t, err.Error(), "some db failure")
}

func TestGetUserByID_Error(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	userID := setup.DefaultTestUser().ID
	expectedErr := gorm.ErrRecordNotFound

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), userID).
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByID, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByID(context.Background(), userID)

	require.Error(t, err)
	require.Equal(t, domain.User{}, userDomain)
}
