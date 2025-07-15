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
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetUserByEmail, commonkeys.Error, expectedErr.Error())

	userDomain, err := suite.UserService.GetUserByEmail(context.Background(), userEmail)

	require.Error(t, err)
	require.Equal(t, domain.User{}, userDomain)
}
