package user_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     "Felipe",
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	suite.UserRepo.EXPECT().
		UpdateUser(suite.Ctx, input.ID, gomock.AssignableToTypeOf(map[string]interface{}{})).
		Return(setup.TestPerfectUser, nil)

	result, err := suite.UserSvc.UpdateUser(suite.Ctx, input)

	assert.NoError(t, err)
	assert.Equal(t, setup.TestPerfectUser, result)
}

func TestUpdateUser_NoFieldsToUpdate(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{ID: setup.TestPerfectUser.ID}

	result, err := suite.UserSvc.UpdateUser(suite.Ctx, input)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Equal(t, constants.ErrorNoFieldsToUpdate, err.Error())
}

func TestUpdateUserPassword_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"
	hashedPassword := "hashedNewPassword"
	expectedToken := "newToken"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepo.EXPECT().
		UpdateUser(
			suite.Ctx,
			input.ID,
			gomock.AssignableToTypeOf(map[string]interface{}{}),
		).Return(setup.TestPerfectUser, nil)

	suite.TokenSvc.EXPECT().
		Create(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return(expectedToken, nil)

	suite.TokenSvc.EXPECT().
		Save(suite.Ctx, domain.TokenDomain{UserID: input.ID, Token: expectedToken}).
		Return(nil)

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.NoError(t, err)
	assert.Equal(t, setup.TestPerfectUser, result)
	assert.Equal(t, expectedToken, token)
}

func TestUpdateUserPassword_ErrorToGetUserByID(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToGetUserByID))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToGetUserByID, err.Error())
}

func TestUpdateUserPassword_ErrorToCompareHashAndPassword(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(errors.New(constants.ErrorToCompareHashAndPassword))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToCompareHashAndPassword, err.Error())
}

func TestUpdateUserPassword_ErrorToHashPassword(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(newPassword).
		Return("", errors.New(constants.ErrorToHashPassword))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestUpdateUserPassword_ErrorToUpdatePassword(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	hashedPassword := "hashedNewPassword"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepo.EXPECT().
		UpdateUser(
			suite.Ctx,
			input.ID,
			gomock.AssignableToTypeOf(map[string]interface{}{}),
		).Return(domain.UserDomain{}, errors.New(constants.ErrorToUpdatePassword))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToUpdatePassword, err.Error())
}

func TestUpdateUserPassword_ErrorToCreateToken(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	hashedPassword := "hashedNewPassword"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepo.EXPECT().
		UpdateUser(
			suite.Ctx,
			input.ID,
			gomock.AssignableToTypeOf(map[string]interface{}{}),
		).Return(setup.TestPerfectUser, nil)

	suite.TokenSvc.EXPECT().
		Create(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return("", errors.New(constants.ErrorToCreateToken))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToCreateToken, err.Error())
}

func TestUpdateUserPassword_ErrorToSaveToken(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		ID:       setup.TestPerfectUser.ID,
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}

	oldPassword := "oldPassword"
	newPassword := "newPassword"

	hashedPassword := "hashedNewPassword"
	expectedToken := "newToken"

	suite.UserRepo.EXPECT().
		GetUserByID(suite.Ctx, input.ID).
		Return(setup.TestPerfectUser, nil)

	suite.PasswordSvc.EXPECT().
		ComparePasswords(setup.TestPerfectUser.Password, oldPassword).
		Return(nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(newPassword).
		Return(hashedPassword, nil)

	suite.UserRepo.EXPECT().
		UpdateUser(
			suite.Ctx,
			input.ID,
			gomock.AssignableToTypeOf(map[string]interface{}{}),
		).Return(setup.TestPerfectUser, nil)

	suite.TokenSvc.EXPECT().
		Create(suite.Ctx, domain.TokenDomain{UserID: input.ID}).
		Return(expectedToken, nil)

	suite.TokenSvc.EXPECT().
		Save(suite.Ctx, domain.TokenDomain{UserID: input.ID, Token: expectedToken}).
		Return(errors.New(constants.ErrorToSaveToken))

	result, token, err := suite.UserSvc.UpdateUserPassword(suite.Ctx, input, oldPassword, newPassword)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, result)
	assert.Empty(t, token)
	assert.Equal(t, constants.ErrorToSaveToken, err.Error())
}
