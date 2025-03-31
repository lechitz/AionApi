package user_test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"testing"
)

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	userID := uint64(1)

	userRepo.EXPECT().SoftDeleteUser(ctx, userID).Return(nil)
	tokenSvc.EXPECT().Delete(ctx, domain.TokenDomain{UserID: userID}).Return(nil)

	err := userSvc.SoftDeleteUser(ctx, userID)

	assert.NoError(t, err)
}

func TestDeleteUser_ErrorToSoftDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	userID := uint64(1)
	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)

	userRepo.EXPECT().SoftDeleteUser(ctx, userID).Return(expectedErr)

	err := userSvc.SoftDeleteUser(ctx, userID)

	assert.EqualError(t, err, expectedErr.Error())
}

func TestDeleteUser_ErrorToDeleteToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zaptest.NewLogger(t).Sugar()

	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)

	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	ctx := domain.ContextControl{}
	userID := uint64(1)
	expectedErr := errors.New(constants.ErrorToDeleteToken)

	userRepo.EXPECT().SoftDeleteUser(ctx, userID).Return(nil)
	tokenSvc.EXPECT().Delete(ctx, domain.TokenDomain{UserID: userID}).Return(expectedErr)

	err := userSvc.SoftDeleteUser(ctx, userID)

	assert.EqualError(t, err, expectedErr.Error())
}
