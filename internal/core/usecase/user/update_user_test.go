// Package user_test contains tests for user use cases.
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
)

func TestUpdateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{
		ID:       setup.DefaultTestUser().ID,
		Name:     "Felipe",
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	expected := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessUserUpdated, commonkeys.UserID, gomock.Any(), commonkeys.UserUpdatedFields, gomock.Any())

	result, err := suite.UserService.UpdateUser(context.Background(), input)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUpdateUser_UpdateOnlyUsername(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{
		ID:       setup.DefaultTestUser().ID,
		Username: "new_username",
	}
	expected := setup.DefaultTestUser()
	expected.Username = "new_username"

	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)
	suite.Logger.EXPECT().InfowCtx(gomock.Any(), constants.SuccessUserUpdated, commonkeys.UserID, gomock.Any(), commonkeys.UserUpdatedFields, gomock.Any())

	result, err := suite.UserService.UpdateUser(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUpdateUser_UpdateOnlyEmail(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{
		ID:    setup.DefaultTestUser().ID,
		Email: "new@email.com",
	}
	expected := setup.DefaultTestUser()
	expected.Email = "new@email.com"

	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)
	suite.Logger.EXPECT().InfowCtx(gomock.Any(), constants.SuccessUserUpdated, commonkeys.UserID, gomock.Any(), commonkeys.UserUpdatedFields, gomock.Any())

	result, err := suite.UserService.UpdateUser(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUpdateUser_NoFieldsToUpdate(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{ID: setup.DefaultTestUser().ID}

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorNoFieldsToUpdate, commonkeys.Error, constants.ErrorNoFieldsToUpdate)

	result, err := suite.UserService.UpdateUser(context.Background(), input)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Equal(t, constants.ErrorNoFieldsToUpdate, err.Error())
}

func TestUpdateUser_ErrorToUpdateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{
		ID:       setup.DefaultTestUser().ID,
		Name:     "Felipe",
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	expectedErr := errors.New(constants.ErrorToUpdateUser)

	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToUpdateUser, commonkeys.Error, expectedErr.Error())

	result, err := suite.UserService.UpdateUser(context.Background(), input)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Contains(t, err.Error(), constants.ErrorToUpdateUser)
}
