// Package setup provides test suite builders and common test helpers for unit tests.
package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/auth/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// AuthServiceTestSuite groups all mocked dependencies and the system under test (AuthService),
// making Auth-related unit tests simpler and more maintainable.
type AuthServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mocks.MockContextLogger
	UserRepository *mocks.MockUserRepository
	Hasher         *mocks.MockHasher
	TokenStore     *mocks.MockAuthStore
	TokenProvider  *mocks.MockAuthProvider
	AuthService    *usecase.Service
	Ctx            context.Context
}

// AuthServiceTest initializes and returns an AuthServiceTestSuite with the correct mocked
// output ports (UserRepository, TokenStore, Hasher, TokenProvider, ContextLogger).
// Use this helper to bootstrap each test and ensure proper teardown via Ctrl.Finish().
func AuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	userRepo := mocks.NewMockUserRepository(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenStore := mocks.NewMockAuthStore(ctrl)
	tokenProvider := mocks.NewMockAuthProvider(ctrl)
	log := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger (no-ops).
	ExpectLoggerDefaultBehavior(log)

	authService := usecase.NewService(
		userRepo,
		tokenStore,
		tokenProvider,
		hasher,
		log,
	)

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         log,
		UserRepository: userRepo,
		Hasher:         hasher,
		TokenStore:     tokenStore,
		TokenProvider:  tokenProvider,
		AuthService:    authService,
		Ctx:            t.Context(),
	}
}
