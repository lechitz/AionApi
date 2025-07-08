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
	UserRetriever  *mocks.MockUserRetriever
	PasswordHasher *mocks.MockHasherStore
	TokenService   *mocks.MockTokenUsecase
	AuthService    *auth.Service
	Ctx            context.Context
}

// AuthServiceTest initializes and returns a test suite with mock dependencies for testing authentication services.
func AuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockUserRetriever := mocks.NewMockUserRetriever(ctrl)
	mockSecurityStore := mocks.NewMockHasherStore(ctrl)
	mockTokenUsecase := mocks.NewMockTokenUsecase(ctrl)
	mockLog := mocks.NewMockLogger(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	authService := auth.NewAuthService(
		mockUserRetriever,
		mockTokenUsecase,
		mockSecurityStore,
		mockLog,
	)

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         mockLog,
		UserRetriever:  mockUserRetriever,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUsecase,
		AuthService:    authService,
		Ctx:            t.Context(),
	}
}
