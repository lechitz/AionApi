package setup

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	mockSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	mockUser "github.com/lechitz/AionApi/tests/mocks/user"
)

// UserServiceTestSuite is a test suite for testing the UserService and its dependencies.
type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mockLogger.MockLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *mockSecurity.MockSecurityStore
	TokenService   *mockToken.MockTokenUsecase
	UserService    *user.Service
	Ctx            context.Context
}

// UserServiceTest initializes and returns a UserServiceTestSuite with mocked dependencies to facilitate unit testing of the user service.
func UserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserRepo := mockUser.NewMockUserStore(ctrl)
	mockSecurityStore := mockSecurity.NewMockSecurityStore(ctrl)
	mockTokenUsecase := mockToken.NewMockTokenUsecase(ctrl)
	mockLog := mockLogger.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	userService := user.NewUserService(mockUserRepo, mockTokenUsecase, mockSecurityStore, mockLog)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRepository: mockUserRepo,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUsecase,
		UserService:    userService,
		Ctx:            t.Context(),
	}
}

// DefaultTestUser is a predefined instance of domain.UserDomain used for testing purposes, representing a perfect/valid user with complete and valid fields.
func DefaultTestUser() domain.UserDomain {
	return domain.UserDomain{
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
