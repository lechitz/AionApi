package user_test

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestCreateUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	input := domain.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@example.com ",
	}
	password := "123"

	normalized := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
		Password: "hashed123",
	}

	userRepo.EXPECT().GetUserByUsername(ctx, "lechitz").Return(domain.UserDomain{}, nil)
	userRepo.EXPECT().GetUserByEmail(ctx, "lechitz@example.com").Return(domain.UserDomain{}, nil)
	passwordSvc.EXPECT().HashPassword(password).Return("hashed123", nil)
	userRepo.EXPECT().CreateUser(ctx, normalized).Return(normalized, nil)

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.NoError(t, err)
	assert.Equal(t, normalized, createdUser)
}

func TestCreateUser_ErrorToGetUserByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	input := domain.UserDomain{
		Name:     "Felipe",
		Username: " lechitz",
		Email:    "lechitz@example.com",
	}
	password := "123"

	userRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := userSvc.CreateUser(ctx, input, password)
	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.UsernameIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToGetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	input := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
	}
	password := "123"

	userRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(domain.UserDomain{}, nil)

	userRepo.EXPECT().
		GetUserByEmail(ctx, "lechitz@example.com").
		Return(domain.UserDomain{ID: 1}, nil)

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.EmailIsAlreadyInUse, err.Error())
}

func TestCreateUser_ErrorToHashPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}

	input := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
	}
	password := "123"

	userRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(domain.UserDomain{}, nil)

	userRepo.EXPECT().
		GetUserByEmail(ctx, "lechitz@example.com").
		Return(domain.UserDomain{}, nil)

	passwordSvc.EXPECT().
		HashPassword(password).
		Return("", errors.New(constants.ErrorToHashPassword))

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToHashPassword, err.Error())
}

func TestCreateUser_ErrorToCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}

	input := domain.UserDomain{
		Name:     "Felipe",
		Username: "lechitz",
		Email:    "lechitz@example.com",
	}
	password := "123"

	userRepo.EXPECT().
		GetUserByUsername(ctx, "lechitz").
		Return(domain.UserDomain{}, nil)

	userRepo.EXPECT().
		GetUserByEmail(ctx, "lechitz@example.com").
		Return(domain.UserDomain{}, nil)

	passwordSvc.EXPECT().
		HashPassword(password).
		Return("hashed123", nil)

	expectedUser := input
	expectedUser.Password = "hashed123"

	userRepo.EXPECT().
		CreateUser(ctx, expectedUser).
		Return(domain.UserDomain{}, errors.New(constants.ErrorToCreateUser))

	createdUser, err := userSvc.CreateUser(ctx, input, password)

	assert.Error(t, err)
	assert.Equal(t, domain.UserDomain{}, createdUser)
	assert.Equal(t, constants.ErrorToCreateUser, err.Error())
}
