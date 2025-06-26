package user_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := entity.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@example.com ",
	}
	password := setup.DefaultTestUser().Password

	normalized := entity.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
		Password: "hashed123",
	}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(entity.UserDomain{}, nil)
	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, "lechitz@example.com").
		Return(entity.UserDomain{}, nil)
	suite.PasswordHasher.EXPECT().
		HashPassword(password).Return("hashed123", nil)
	suite.UserRepository.EXPECT().
		CreateUser(suite.Ctx, normalized).Return(normalized, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	require.NoError(t, err)
	require.Equal(t, normalized, createdUser)
}

func TestCreateUser_ErrorToGetUserByUsername(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := entity.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.DefaultTestUser().Username).
		Return(entity.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	require.Error(t, err)
	require.Equal(t, entity.UserDomain{}, createdUser)
	require.Equal(t, constants.UsernameIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToGetUserByEmail(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := entity.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.DefaultTestUser().Username).
		Return(entity.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.DefaultTestUser().Email).
		Return(entity.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	require.Error(t, err)
	require.Equal(t, entity.UserDomain{}, createdUser)
	require.Equal(t, constants.EmailIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := entity.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.DefaultTestUser().Username).
		Return(entity.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.DefaultTestUser().Email).
		Return(entity.UserDomain{}, nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("", errors.New(constants.ErrorToHashPassword))

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	require.Error(t, err)
	require.Equal(t, entity.UserDomain{}, createdUser)
	require.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := entity.UserDomain{
		Name:     setup.DefaultTestUser().Name,
		Username: setup.DefaultTestUser().Username,
		Email:    setup.DefaultTestUser().Email,
	}
	password := setup.DefaultTestUser().Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.DefaultTestUser().Username).
		Return(entity.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.DefaultTestUser().Email).
		Return(entity.UserDomain{}, nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("hashed123", nil)

	expectedUser := input
	expectedUser.Password = "hashed123"

	suite.UserRepository.EXPECT().
		CreateUser(suite.Ctx, expectedUser).
		Return(entity.UserDomain{}, errors.New(constants.ErrorToCreateUser))

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	require.Error(t, err)
	require.Equal(t, entity.UserDomain{}, createdUser)
	require.Equal(t, constants.ErrorToCreateUser, err.Error())
}
