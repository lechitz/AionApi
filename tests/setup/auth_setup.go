package setup

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	mockSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	mockUser "github.com/lechitz/AionApi/tests/mocks/user"
)

type AuthServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mockLogger.MockLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *mockSecurity.MockSecurityStore
	TokenService   *mockToken.MockTokenUsecase
	AuthService    *auth.AuthService
	Ctx            context.Context
}

func SetupAuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserRepo := mockUser.NewMockUserStore(ctrl)
	mockSecurityStore := mockSecurity.NewMockSecurityStore(ctrl)
	mockTokenUseCase := mockToken.NewMockTokenUsecase(ctrl)
	mockLog := mockLogger.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	authService := auth.NewAuthService(mockUserRepo, mockTokenUseCase, mockSecurityStore, mockLog, "supersecretkey")

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRepository: mockUserRepo,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUseCase,
		AuthService:    authService,
		Ctx:            context.Background(),
	}
}
