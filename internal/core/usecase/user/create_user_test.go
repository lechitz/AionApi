// Package user_test contains tests for user use cases.
package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@example.com ",
	}
	password := setup.DefaultTestUser().Password

	normalized := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
		Password: "hashed123",
	}

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), "lechitz").
		Return(domain.UserDomain{}, nil)
	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), "lechitz@example.com").
		Return(domain.UserDomain{}, nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(password).Return("hashed123", nil)
	suite.UserRepository.EXPECT().
		CreateUser(gomock.Any(), normalized).Return(normalized, nil)

	suite.Logger.EXPECT().InfowCtx(gomock.Any(), constants.SuccessUserCreated, commonkeys.UserID, gomock.Any())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.NoError(t, err)
	require.Equal(t, normalized, createdUser)
}

func TestCreateUser_UsernameAlreadyExists(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
	require.Equal(t, "validation error on username: "+sharederrors.ErrUsernameInUse, err.Error())
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), setup.DefaultTestUser().Email).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
	require.Equal(t, "validation error on email: "+sharederrors.ErrEmailInUse, err.Error())
}

func TestCreateUser_ErrorToGetUserByUsername(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	expectedErr := errors.New(constants.DBErrorCheckingUsername)
	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.DBErrorCheckingUsername, commonkeys.Error, expectedErr.Error())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
}

func TestCreateUser_ErrorToGetUserByEmail(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{}, nil)

	expectedErr := errors.New(constants.DBErrorCheckingEmail)
	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), setup.DefaultTestUser().Email).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.DBErrorCheckingEmail, commonkeys.Error, expectedErr.Error())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), setup.DefaultTestUser().Email).
		Return(domain.UserDomain{}, nil)

	expectedErr := errors.New(constants.ErrorToHashPassword)
	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("", expectedErr)

	suite.Logger.EXPECT().ErrorwCtx(gomock.Any(), constants.ErrorToHashPassword, commonkeys.Error, expectedErr.Error())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(gomock.Any(), setup.DefaultTestUser().Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(gomock.Any(), setup.DefaultTestUser().Email).
		Return(domain.UserDomain{}, nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("hashed123", nil)

	expectedUser := input
	expectedUser.Password = "hashed123"

	expectedErr := errors.New(constants.ErrorToCreateUser)
	suite.UserRepository.EXPECT().
		CreateUser(gomock.Any(), expectedUser).
		Return(domain.UserDomain{}, expectedErr)

	suite.Logger.EXPECT().ErrorwCtx(gomock.Any(), constants.ErrorToCreateUser, commonkeys.Error, expectedErr.Error())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, expectedErr, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
}

func TestCreateUser_ValidationErrorRequiredFields(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     "",
		Username: "",
		Email:    "",
	}
	password := ""

	suite.Logger.EXPECT().
		ErrorwCtx(gomock.Any(), constants.ErrorToValidateCreateUser, commonkeys.Error, gomock.Any())

	createdUser, err := suite.UserService.CreateUser(context.Background(), input, password)

	require.Error(t, err)
	require.Equal(t, domain.UserDomain{}, createdUser)
	require.Contains(t, err.Error(), "validation error")
	require.Contains(t, err.Error(), "is required")
}
