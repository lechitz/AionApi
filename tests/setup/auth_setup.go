package setup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	mockSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	mockUser "github.com/lechitz/AionApi/tests/mocks/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type AuthServiceTestSuite struct {
	Ctrl           *gomock.Controller
	LoggerSugar    *zap.SugaredLogger
	UserRepository *mockUser.MockUserStore
	PasswordHasher *mockSecurity.MockSecurityStore
	TokenService   *mockToken.MockTokenUsecase
	AuthService    *auth.AuthService
	Ctx            domain.ContextControl
}

func SetupAuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)
	logger := zaptest.NewLogger(t).Sugar()

	mockUserRepo := mockUser.NewMockUserStore(ctrl)
	mockSecurityStore := mockSecurity.NewMockSecurityStore(ctrl)
	mockTokenUsecase := mockToken.NewMockTokenUsecase(ctrl)

	authService := auth.NewAuthService(mockUserRepo, mockTokenUsecase, mockSecurityStore, logger, "supersecretkey")

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		LoggerSugar:    logger,
		UserRepository: mockUserRepo,
		PasswordHasher: mockSecurityStore,
		TokenService:   mockTokenUsecase,
		AuthService:    authService,
		Ctx:            domain.ContextControl{},
	}
}
