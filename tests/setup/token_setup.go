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
	UserCache      *mocks.MockUserCache
	RoleCache      *mocks.MockRoleCache
}

// TokenServiceTest initializes and returns a TokenServiceTestSuite with mocked output ports.
func TokenServiceTest(t *testing.T) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	userRepository := mocks.NewMockUserRepository(ctrl)
	userCache := mocks.NewMockUserCache(ctrl)
	logger := mocks.NewMockContextLogger(ctrl)
	hasher := mocks.NewMockHasher(ctrl)
	authStore := mocks.NewMockAuthStore(ctrl)
	authProvider := mocks.NewMockAuthProvider(ctrl)
	rolesReader := mocks.NewMockRolesReader(ctrl)
	roleCache := mocks.NewMockRoleCache(ctrl)

	ExpectLoggerDefaultBehavior(logger)
	roleCache.EXPECT().GetRoles(gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	roleCache.EXPECT().SaveRoles(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	svc := usecase.NewService(rolesReader, roleCache, userRepository, userCache, authStore, authProvider, hasher, logger)

	return &TokenServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         logger,
		TokenStore:     authStore,
		AuthProvider:   authProvider,
		UserRepository: userRepository,
		UserCache:      userCache,
		RoleCache:      roleCache,
		TokenService:   svc,
		Ctx:            t.Context(),
	}
}
