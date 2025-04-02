package setup

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	mocksSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mocksToken "github.com/lechitz/AionApi/tests/mocks/token"
	mocksDB "github.com/lechitz/AionApi/tests/mocks/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

type AuthServiceTestSuite struct {
	Ctrl           *gomock.Controller
	Logger         *zap.SugaredLogger
	UserRepo       *mocksDB.MockRepository
	TokenService   *mocksToken.MockStore
	PasswordHasher *mocksSecurity.MockHasher
	AuthService    *auth.AuthService
	Ctx            domain.ContextControl
}

func SetupAuthServiceTest(t *testing.T) *AuthServiceTestSuite {
	ctrl := gomock.NewController(t)

	logger := zaptest.NewLogger(t).Sugar()
	userRepo := mocksDB.NewMockRepository(ctrl)
	tokenService := mocksToken.NewMockStore(ctrl)
	passwordHasher := mocksSecurity.NewMockHasher(ctrl)

	authService := &auth.AuthService{
		UserRetriever:  userRepo,
		TokenService:   tokenService,
		PasswordHasher: passwordHasher,
		LoggerSugar:    logger,
	}

	return &AuthServiceTestSuite{
		Ctrl:           ctrl,
		Logger:         logger,
		UserRepo:       userRepo,
		TokenService:   tokenService,
		PasswordHasher: passwordHasher,
		AuthService:    authService,
		Ctx:            domain.ContextControl{},
	}
}

var (
	TestPerfectAuthUser = domain.UserDomain{
		ID:        1,
		Name:      "Auth Test User",
		Username:  "user",
		Email:     "user@example.com",
		Password:  "supersecure123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}
)
