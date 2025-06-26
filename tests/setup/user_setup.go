package setup

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/user"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	securitymocks "github.com/lechitz/AionApi/tests/mocks/security"
	tokenmocks "github.com/lechitz/AionApi/tests/mocks/token"
	mockUser "github.com/lechitz/AionApi/tests/mocks/user"
	"go.uber.org/mock/gomock"
)

// UserServiceTestSuite is a test suite for testing the UserService and its dependencies.
type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mockLogger.MockLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *securitymocks.MockSecurityStore
	TokenService   *tokenmocks.MockTokenUsecase
	UserService    *user.Service
	Ctx            context.Context
}

// UserServiceTest initializes and returns a UserServiceTestSuite with mocked dependencies to facilitate unit testing of the user service.
func UserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserRepo := mockUser.NewMockUserStore(ctrl)
	mockSecurityStore := securitymocks.NewMockSecurityStore(ctrl)
	mockTokenUseCase := tokenmocks.NewMockTokenUsecase(ctrl)
	mockLog := mockLogger.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	userService := user.NewUserService(mockUserRepo, mockTokenUseCase, mockSecurityStore, mockLog)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRepository: mockUserRepo,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUseCase,
		UserService:    userService,
		Ctx:            t.Context(),
	}
}

// DefaultTestUser is a predefined instance of domain.UserDomain used for testing purposes, representing a perfect/valid user with complete and valid fields.
func DefaultTestUser() entity.UserDomain {
	return entity.UserDomain{
		ID:        1,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "user@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}
}
