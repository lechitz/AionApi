package user

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestUpdateUserPassword_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	expectedUser := input
	expectedUser.Password = "password123"
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"
	expectedToken := "newToken"

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(expectedUser, nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(expectedUser.Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)
	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expectedUser, nil)
	suite.TokenService.EXPECT().
		CreateToken(gomock.Any(), gomock.AssignableToTypeOf(domain.Token{})).
		Return(expectedToken, nil)
	suite.Logger.EXPECT().
		InfowCtx(gomock.Any(), constants.SuccessPasswordUpdated, commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.NoError(t, err)
	require.Equal(t, expectedUser, result)
	require.Equal(t, expectedToken, token)
}

func TestUpdateUserPassword_ErrorToGetUserByID(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"

	expectedErr := errors.New(constants.ErrorToGetSelf)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToGetSelf, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToGetSelf)
}

func TestUpdateUserPassword_ErrorToCompareHashAndPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"

	expectedErr := errors.New(constants.ErrorToCompareHashAndPassword)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToCompareHashAndPassword, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToCompareHashAndPassword)
}

func TestUpdateUserPassword_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"

	expectedErr := errors.New(constants.ErrorToHashPassword)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return("", expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToHashPassword, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToHashPassword)
}

func TestUpdateUserPassword_ErrorToUpdatePassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"

	expectedErr := errors.New(constants.ErrorToUpdatePassword)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)
	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(domain.User{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToUpdatePassword, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToUpdatePassword)
}

func TestUpdateUserPassword_ErrorToCreateToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"

	expectedErr := errors.New(constants.ErrorToCreateToken)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)
	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.DefaultTestUser(), nil)
	suite.TokenService.EXPECT().
		CreateToken(gomock.Any(), domain.Token{Key: input.ID}).
		Return("", expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToCreateToken, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToCreateToken)
}

func TestUpdateUserPassword_ErrorToSaveToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.User{
		ID:       setup.DefaultTestUser().ID,
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"

	expectedErr := errors.New(constants.ErrorToSaveToken)
	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)
	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.DefaultTestUser(), nil)
	suite.TokenService.EXPECT().
		CreateToken(gomock.Any(), domain.Token{Key: input.ID}).
		Return("", expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToCreateToken, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), constants.ErrorToSaveToken)
}

func TestUpdateUserPassword_UnknownTokenServiceError(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"
	expectedErr := errors.New("random unknown error")

	suite.UserRepository.EXPECT().
		GetUserByID(gomock.Any(), input.ID).
		Return(setup.DefaultTestUser(), nil)
	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)
	suite.UserRepository.EXPECT().
		UpdateUser(gomock.Any(), input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.DefaultTestUser(), nil)
	suite.TokenService.EXPECT().
		CreateToken(gomock.Any(), domain.Token{Key: input.ID}).
		Return("", expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToCreateToken, commonkeys.Error, expectedErr.Error(), commonkeys.UserID, gomock.Any())

	result, token, err := suite.UserService.UpdatePassword(
		context.Background(),
		input,
		oldPassword,
		newPassword,
	)
	require.Error(t, err)
	require.Equal(t, domain.User{}, result)
	require.Empty(t, token)
	require.Contains(t, err.Error(), expectedErr.Error())
}
