package setup

import (
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/tests/mocks"

	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"go.uber.org/mock/gomock"
)

// UserServiceTestSuite is a test suite for testing the UserService and its dependencies.
type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mocks.MockLogger
	UserRepository *mocks.MockUserStore
	PasswordHasher *mocks.MockHasherStore
	TokenService   *mocks.MockTokenUsecase
	UserService    *user.Service
}

// UserServiceTest initializes and returns a UserServiceTestSuite with mocked dependencies to facilitate unit testing of the user service.
func UserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserStore := mocks.NewMockUserStore(ctrl)
	mockSecurityStore := mocks.NewMockHasherStore(ctrl)
	mockTokenUsecase := mocks.NewMockTokenUsecase(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	userService := user.NewService(
		mockUserStore,
		mockTokenUsecase,
		mockSecurityStore,
		mockLog,
	)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRepository: mockUserStore,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUsecase,
		UserService:    userService,
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
