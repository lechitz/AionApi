package setup

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	mockSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	mockUser "github.com/lechitz/AionApi/tests/mocks/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	LoggerSugar    *zap.SugaredLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *mockSecurity.MockSecurityStore
	TokenService   *mockToken.MockTokenUsecase
	UserService    *user.UserService
	Ctx            domain.ContextControl
}

func SetupUserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)
	logger := zaptest.NewLogger(t).Sugar()

	mockUserRepo := mockUser.NewMockUserStore(ctrl)
	mockSecurityStore := mockSecurity.NewMockSecurityStore(ctrl)
	mockTokenUsecase := mockToken.NewMockTokenUsecase(ctrl)

	userService := user.NewUserService(mockUserRepo, mockTokenUsecase, mockSecurityStore, logger)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		LoggerSugar:    logger,
		UserRepository: mockUserRepo,
		TokenService:   mockTokenUsecase,
		PasswordHasher: mockSecurityStore,
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
