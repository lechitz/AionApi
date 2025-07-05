package user_test

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.DefaultTestUser().ID,
		Name:     "Felipe",
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}

	expected := setup.DefaultTestUser()

	suite.UserRepository.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expected, nil)

	result, err := suite.UserService.UpdateUser(suite.Ctx, input)

	require.NoError(t, err)
	require.Equal(t, expected, result)
}

func TestUpdateUser_NoFieldsToUpdate(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{ID: setup.DefaultTestUser().ID}

	result, err := suite.UserService.UpdateUser(suite.Ctx, input)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Equal(t, constants.ErrorNoFieldsToUpdate, err.Error())
}

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
		GetUserByID(suite.Ctx, input.ID).
		Return(expectedUser, nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(expectedUser.Password, oldPassword).
		Return(nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepository.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(expectedUser, nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return(expectedToken, nil)

	suite.TokenService.EXPECT().
		Save(suite.Ctx, domain.TokenDomain{UserID: input.ID, Token: expectedToken}).
		Return(nil)

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
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

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToGetUserByID))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToGetUserByID, err.Error())
}

func TestUpdateUserPassword_ErrorToCompareHashAndPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.DefaultTestUser(), nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(errors.New(constants.ErrorToCompareHashAndPassword))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToCompareHashAndPassword, err.Error())
}

func TestUpdateUserPassword_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.DefaultTestUser(), nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return("", errors.New(constants.ErrorToHashPassword))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestUpdateUserPassword_ErrorToUpdatePassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.DefaultTestUser(), nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepository.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToUpdatePassword))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToUpdatePassword, err.Error())
}

func TestUpdateUserPassword_ErrorToCreateToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := setup.DefaultTestUser()
	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.DefaultTestUser(), nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepository.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.DefaultTestUser(), nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return("", errors.New(constants.ErrorToCreateToken))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToCreateToken, err.Error())
}

func TestUpdateUserPassword_ErrorToSaveToken(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.DefaultTestUser().ID,
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"
	expectedToken := "newToken"

	suite.UserRepository.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.DefaultTestUser(), nil)

	suite.PasswordHasher.EXPECT().
		ValidatePassword(setup.DefaultTestUser().Password, oldPassword).
		Return(nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepository.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.DefaultTestUser(), nil)

	suite.TokenService.EXPECT().
		CreateToken(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return(expectedToken, nil)

	suite.TokenService.EXPECT().
		Save(suite.Ctx, domain.TokenDomain{UserID: input.ID, Token: expectedToken}).
		Return(errors.New(constants.ErrorToSaveToken))

	result, token, err := suite.UserService.UpdateUserPassword(
		suite.Ctx,
		input,
		oldPassword,
		newPassword,
	)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, result)
	require.Empty(t, token)
	require.Equal(t, constants.ErrorToSaveToken, err.Error())
}
