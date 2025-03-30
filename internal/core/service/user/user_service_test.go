package user

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/service/constants"
	"github.com/lechitz/AionApi/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
)

func TestNormalizeUserData(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	service := &UserService{LoggerSugar: logger}

	raw := &domain.UserDomain{
		Name:     "  Felipe  ",
		Username: " lechitz ",
		Email:    "  LECHITZ@AION.COM ",
	}

	normalized := service.NormalizeUserData(raw)

	assert.Equal(t, "Felipe", normalized.Name)
	assert.Equal(t, "lechitz", normalized.Username)
	assert.Equal(t, "lechitz@aion.com", normalized.Email)
}

func TestValidateCreateUserRequired(t *testing.T) {
	service := &UserService{}

	tests := []struct {
		name     string
		user     domain.UserDomain
		password string
		wantErr  string
	}{
		{"missing name", domain.UserDomain{Username: "a", Email: "a@b.com"}, "123", constants.NameIsRequired},
		{"missing username", domain.UserDomain{Name: "a", Email: "a@b.com"}, "123", constants.UsernameIsRequired},
		{"missing email", domain.UserDomain{Name: "a", Username: "a"}, "123", constants.EmailIsRequired},
		{"invalid email", domain.UserDomain{Name: "a", Username: "a", Email: "invalid"}, "123", constants.InvalidEmail},
		{"missing password", domain.UserDomain{Name: "a", Username: "a", Email: "a@b.com"}, "", constants.PasswordIsRequired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateCreateUserRequired(tt.user, tt.password)
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestSoftDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTokenService := mocks.NewMockTokenServiceInterface(ctrl)
	mockPasswordService := mocks.NewMockPasswordManager(ctrl)

	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(42)

	mockUserRepo.EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(nil)

	mockTokenService.EXPECT().
		Delete(ctx, domain.TokenDomain{UserID: userID}).
		Return(nil)

	err := service.SoftDeleteUser(ctx, userID)
	assert.NoError(t, err)
}

func TestSoftDeleteUser_FailureOnUserRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTokenService := mocks.NewMockTokenServiceInterface(ctrl)
	mockPasswordService := mocks.NewMockPasswordManager(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(99)
	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)

	mockUserRepo.EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(expectedErr)

	err := service.SoftDeleteUser(ctx, userID)
	assert.EqualError(t, err, expectedErr.Error())
}

func TestSoftDeleteUser_FailureOnTokenDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTokenService := mocks.NewMockTokenServiceInterface(ctrl)
	mockPasswordService := mocks.NewMockPasswordManager(ctrl)
	logger := zaptest.NewLogger(t).Sugar()

	service := NewUserService(mockUserRepo, mockTokenService, mockPasswordService, logger)

	ctx := domain.ContextControl{}
	userID := uint64(77)
	expectedErr := errors.New(constants.ErrorToSoftDeleteUser)

	mockUserRepo.EXPECT().
		SoftDeleteUser(ctx, userID).
		Return(nil)

	mockTokenService.EXPECT().
		Delete(ctx, domain.TokenDomain{UserID: userID}).
		Return(expectedErr)

	err := service.SoftDeleteUser(ctx, userID)
	assert.EqualError(t, err, expectedErr.Error())
}
