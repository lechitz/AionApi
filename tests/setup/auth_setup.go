package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// AuthServiceTestSuite defines a test suite for AuthService, including mock services and dependencies for testing authentication components.
type AuthServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mocks.MockLogger
	UserRepository *mocks.MockUserStore
	PasswordHasher *mocks.MockSecurityStore
	TokenService   *mocks.MockTokenUsecase
	AuthService    *auth.Service
	Ctx            context.Context
}

// AuthServiceTest initializes and returns a test suite with mock dependencies for testing authentication services.
func AuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserRepo := mocks.NewMockUserStore(ctrl)
	mockSecurityStore := mocks.NewMockSecurityStore(ctrl)
	mockTokenUseCase := mocks.NewMockTokenUsecase(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	authService := auth.NewAuthService(
		mockUserRepo,
		mockTokenUseCase,
		mockSecurityStore,
		mockLog,
		"supersecretkey",
	)

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRepository: mockUserRepo,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUseCase,
		AuthService:    authService,
		Ctx:            t.Context(),
	}
}
