package setup

import (
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

type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mockLogger.MockLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *mockSecurity.MockSecurityStore
	TokenService   *mockToken.MockTokenUsecase
	UserService    *user.UserService
	Ctx            domain.ContextControl
}

func SetupUserServiceTest(t *testing.T) *UserServiceTestSuite {
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
		Ctx:            domain.ContextControl{},
	}
}

var TestPerfectUser = domain.UserDomain{
	ID:        1,
	Name:      "Test User",
	Username:  "testuser",
	Email:     "user@example.com",
	Password:  "password123",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: nil,
}
