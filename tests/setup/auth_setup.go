package setup

//
//import (
//	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
//	"github.com/lechitz/AionApi/internal/core/usecase/token"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"github.com/lechitz/AionApi/internal/core/domain"
//	"github.com/lechitz/AionApi/internal/core/usecase/auth"
//	"github.com/lechitz/AionApi/tests/mocks"
//	"go.uber.org/zap"
//	"go.uber.org/zap/zaptest"
//	"gorm.io/gorm"
//)
//
//type AuthServiceTestSuite struct {
//	Ctrl           *gomock.Controller
//	Logger         *zap.SugaredLogger
//	UserRepo       *mocks.MockUserRepository
//	TokenService   token.TokenService
//	PasswordHasher *mocks.MockSecurityStore
//	CacheConn      *mocks.MockCacheClient
//	AuthService    *auth.AuthService
//	Ctx            domain.ContextControl
//}
//
//func SetupAuthServiceTest(t *testing.T) *AuthServiceTestSuite {
//	ctrl := gomock.NewController(t)
//
//	logger := zaptest.NewLogger(t).Sugar()
//
//	userRepo := mocks.NewMockUserRepository(ctrl)
//
//	client := mocks.NewMockCacheClient(ctrl)
//		tokenRepo := cache.NewTokenRepository(client, logger)
//	tokenService := token.NewTokenService(*tokenRepo, logger, "XYZ1234567890")
//	securityHasher := mocks.NewMockSecurityStore(ctrl)
//
//	authService := auth.NewAuthService(userRepo, *tokenService, securityHasher, logger, "supersecretkey")
//
//	return &AuthServiceTestSuite{
//		Ctrl:           ctrl,
//		Logger:         logger,
//		UserRepo:       userRepo,
//		TokenService:   *tokenService,
//		PasswordHasher: securityHasher,
//		AuthService:    authService,
//		Ctx:            domain.ContextControl{},
//	}
//}
//
//var (
//	TestPerfectAuthUser = domain.UserDomain{
//		ID:        1,
//		Name:      "Auth Test User",
//		Username:  "user",
//		Email:     "user@example.com",
//		Password:  "supersecure123",
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//		DeletedAt: gorm.DeletedAt{},
//	}
//)
