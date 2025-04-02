package user_test

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/tests/setup"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCreateUser_Success(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
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

	suite.UserRepo.EXPECT().GetUserByUsername(suite.Ctx, "lechitz").Return(domain.UserDomain{}, nil)
	suite.UserRepo.EXPECT().GetUserByEmail(suite.Ctx, "lechitz@example.com").Return(domain.UserDomain{}, nil)
	suite.PasswordSvc.EXPECT().HashPassword(password).Return("hashed123", nil)
	suite.UserRepo.EXPECT().CreateUser(suite.Ctx, normalized).Return(normalized, nil)

	createdUser, err := suite.UserSvc.CreateUser(suite.Ctx, input, password)

	assert.NoError(t, err)
	assert.Equal(t, normalized, createdUser)
}

func TestCreateUser_ErrorToGetUserByUsername(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepo.EXPECT().
		GetUserByUsername(suite.Ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := suite.UserSvc.CreateUser(suite.Ctx, input, password)
	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.UsernameIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToGetUserByEmail(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userSvc := user.NewUserService(suite.UserRepo, suite.TokenSvc, suite.PasswordSvc, logger)

	ctx := domain.ContextControl{}
	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepo.EXPECT().
		GetUserByEmail(ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.EmailIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userSvc := user.NewUserService(suite.UserRepo, suite.TokenSvc, suite.PasswordSvc, logger)

	ctx := domain.ContextControl{}

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepo.EXPECT().
		GetUserByEmail(ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{}, nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(password).
		Return("", errors.New(constants.ErrorToHashPassword))

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	suite := setup.SetupUserServiceTest(t)
	defer suite.Ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userSvc := user.NewUserService(suite.UserRepo, suite.TokenSvc, suite.PasswordSvc, logger)

	ctx := domain.ContextControl{}

	input := domain.UserDomain{
		Name:     setup.TestPerfectUser.Name,
		Username: setup.TestPerfectUser.Username,
		Email:    setup.TestPerfectUser.Email,
	}
	password := setup.TestPerfectUser.Password

	suite.UserRepo.EXPECT().
		GetUserByUsername(ctx, setup.TestPerfectUser.Username).
		Return(domain.UserDomain{}, nil)

	suite.UserRepo.EXPECT().
		GetUserByEmail(ctx, setup.TestPerfectUser.Email).
		Return(domain.UserDomain{}, nil)

	suite.PasswordSvc.EXPECT().
		HashPassword(password).
		Return("hashed123", nil)

	expectedUser := input
	expectedUser.Password = "hashed123"

	suite.UserRepo.EXPECT().
		CreateUser(ctx, expectedUser).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToCreateUser))

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToCreateUser, err.Error())
}
