// Package setup provides test suite builders and common test helpers for unit tests.
package setup

import (
	"context"
	"testing"
	"time"

	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/usecase"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// UserServiceTestSuite groups mocked dependencies and the system under test (UserService)
// to keep user-related tests concise and consistent.
type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *mocks.MockContextLogger
	UserRepository *mocks.MockUserRepository
	TokenStore     *mocks.MockAuthStore
	TokenProvider  *mocks.MockAuthProvider
	Hasher         *mocks.MockHasher
	UserService    *usecase.Service
	Ctx            context.Context
}

// UserServiceTest initializes and returns a UserServiceTestSuite using mocked output ports.
// Use this helper to bootstrap each test and ensure proper teardown via Ctrl.Finish().
func UserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	userRepo := mocks.NewMockUserRepository(ctrl)
	tokenStore := mocks.NewMockAuthStore(ctrl)
	tokenProvider := mocks.NewMockAuthProvider(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	log := mocks.NewMockContextLogger(ctrl)

	// Set default, non-intrusive expectations for the logger (no-ops).
	ExpectLoggerDefaultBehavior(log)

	svc := usecase.NewService(
		userRepo,
		tokenStore,
		tokenProvider,
		hasher,
		log,
	)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         log,
		UserRepository: userRepo,
		TokenStore:     tokenStore,
		TokenProvider:  tokenProvider,
		Hasher:         hasher,
		UserService:    svc,
		Ctx:            t.Context(),
	}
}

// DefaultTestUser returns a valid domain.User commonly used in unit tests.
func DefaultTestUser() domain.User {
	return domain.User{
		ID:        1,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "user@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
