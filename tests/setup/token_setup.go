package setup

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
)

// TokenServiceTestSuite groups mocked dependencies and the SUT (TokenService)
// to simplify Token-related unit tests.
type TokenServiceTestSuite struct {
	Ctx            context.Context
	TokenService   *usecase.Service
	Ctrl           *gomock.Controller
	Logger         *mocks.MockContextLogger
	TokenStore     *mocks.MockAuthStore
	AuthProvider   *mocks.MockAuthProvider
	UserRepository *mocks.MockUserRepository
}

// TokenServiceTest initializes and returns a TokenServiceTestSuite with mocked output ports.
func TokenServiceTest(t *testing.T) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	userRepository := mocks.NewMockUserRepository(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	authStore := mocks.NewMockAuthStore(ctrl)
	authProvider := mocks.NewMockAuthProvider(ctrl)

	ExpectLoggerDefaultBehavior(logger)

	svc := usecase.NewService(userRepository, authStore, authProvider, hasher, logger)

	return &TokenServiceTestSuite{
		Ctrl:         ctrl,
		Logger:       logger,
		TokenStore:   authStore,
		AuthProvider: authProvider,
		TokenService: svc,
		Ctx:          t.Context(),
	}
}
