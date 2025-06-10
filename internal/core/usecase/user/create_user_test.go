package user_test

import (
	"errors"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/tests/setup"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@example.com ",
	}
	password := setup.TestPerfectUser.Password

	normalized := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
		Password: "hashed123",
	}

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, "lechitz").
		Return(domain.UserDomain{}, nil)
	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, "lechitz@example.com").
		Return(domain.UserDomain{}, nil)
	suite.PasswordHasher.EXPECT().HashPassword(password).Return("hashed123", nil)
	suite.UserRepository.EXPECT().CreateUser(suite.Ctx, normalized).Return(normalized, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	assert.NoError(t, err)
	assert.Equal(t, normalized, createdUser)
}

func TestCreateUser_ErrorToGetUserByUsername(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.UsernameIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToGetUserByEmail(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.EmailIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{}, nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("", errors.New(constants.ErrorToHashPassword))

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	suite := setup.UserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepository.EXPECT().
		GetUserByUsername(suite.Ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepository.EXPECT().
		GetUserByEmail(suite.Ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{}, nil)

	suite.PasswordHasher.EXPECT().
		HashPassword(password).
		Return("hashed123", nil)

	expectedUser := input
	expectedUser.Password = "hashed123"

	suite.UserRepository.EXPECT().
		CreateUser(suite.Ctx, expectedUser).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToCreateUser))

	createdUser, err := suite.UserService.CreateUser(suite.Ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToCreateUser, err.Error())
}
