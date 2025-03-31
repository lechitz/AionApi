package setup

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/zap/zaptest"
)

type UserServiceTestSuite struct {
	Ctrl        *gomock.Controller
	Logger      *zap.SugaredLogger
	UserRepo    *mocks.MockUserRepository
	PasswordSvc *mocks.MockPasswordManager
	TokenSvc    *mocks.MockTokenServiceInterface
	UserSvc     *user.UserService
	Ctx         domain.ContextControl
}

func SetupUserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	logger := zaptest.NewLogger(t).Sugar()
	userRepo := mocks.NewMockUserRepository(ctrl)
	passwordSvc := mocks.NewMockPasswordManager(ctrl)
	tokenSvc := mocks.NewMockTokenServiceInterface(ctrl)
	userSvc := user.NewUserService(userRepo, tokenSvc, passwordSvc, logger)

	return &UserServiceTestSuite{
		Ctrl:        ctrl,
		Logger:      logger,
		UserRepo:    userRepo,
		PasswordSvc: passwordSvc,
		TokenSvc:    tokenSvc,
		UserSvc:     userSvc,
		Ctx:         domain.ContextControl{},
	}
}

var (
	TestPerfectUser = domain.UserDomain{
		ID:        1,
		Name:      "Test User",
		Username:  "testuser",
		Email:     "user@example.com",
		Password:  "password123",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: gorm.DeletedAt{},
	}
)
