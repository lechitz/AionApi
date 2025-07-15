// Package user_test contains tests for user use cases.
package user_test

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

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
