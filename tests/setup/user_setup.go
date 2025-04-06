package setup

import (
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

type UserServiceTestSuite struct {
	Ctrl           *gomock.Controller
	UserRepository *mocks.MockUserRepository
	HasherStore    *mocks.MockSecurityStore
	TokenService   *mocks.MockTokenService
	UserSvc        *user.UserService
	LoggerSugar    *zap.SugaredLogger
	Ctx            domain.ContextControl
}

func SetupUserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	logger := zaptest.NewLogger(t).Sugar()

	mockUserRepository := mocks.NewMockUserRepository(ctrl)

	tokenService := mocks.NewMockTokenService(ctrl)
	mockHasherStore := mocks.NewMockSecurityStore(ctrl)

	ctx := domain.ContextControl{}

	userSvc := user.NewUserService(mockUserRepository, token.TokenService{}, mockHasherStore, logger)

	return &UserServiceTestSuite{
		Ctrl:           ctrl,
		UserRepository: mockUserRepository,
		HasherStore:    mockHasherStore,
		TokenService:   tokenService,
		UserSvc:        userSvc,
		LoggerSugar:    logger,
		Ctx:            ctx,
	}
}

var TestPerfectUser = domain.UserDomain{
	ID:        1,
	Name:      "Test User",
	Username:  "testuser",
	Email:     "user@example.com",
	Password:  "password123",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
	DeletedAt: gorm.DeletedAt{},
}
