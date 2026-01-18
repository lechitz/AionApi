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
	RolesReader    *mocks.MockRolesReader
	UserRepository *mocks.MockUserRepository
	UserCache      *mocks.MockUserCache
	Hasher         *mocks.MockHasher
	TokenStore     *mocks.MockAuthStore
	TokenProvider  *mocks.MockAuthProvider
	AuthService    *usecase.Service
	Ctx            context.Context
}

// AuthServiceTest initializes and returns an AuthServiceTestSuite with the correct mocked
// output ports (UserRepository, UserCache, TokenStore, Hasher, TokenProvider, ContextLogger).
// Use this helper to bootstrap each test and ensure proper teardown via Ctrl.Finish().
func AuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	userRepo := mocks.NewMockUserRepository(ctrl)
	userCache := mocks.NewMockUserCache(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	tokenStore := mocks.NewMockAuthStore(ctrl)
	tokenProvider := mocks.NewMockAuthProvider(ctrl)
	rolesReader := mocks.NewMockRolesReader(ctrl)
	log := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger (no-ops).
	ExpectLoggerDefaultBehavior(log)

	authService := usecase.NewService(
		rolesReader,
		userRepo,
		userCache,
		tokenStore,
		tokenProvider,
		hasher,
		log,
	)

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         log,
		RolesReader:    rolesReader,
		UserRepository: userRepo,
		UserCache:      userCache,
		Hasher:         hasher,
		TokenStore:     tokenStore,
		TokenProvider:  tokenProvider,
		AuthService:    authService,
		Ctx:            t.Context(),
	}
}
