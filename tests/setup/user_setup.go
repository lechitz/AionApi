package setup

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	mocksSecurity "github.com/lechitz/AionApi/tests/mocks/security"
	mocksToken "github.com/lechitz/AionApi/tests/mocks/token"
	mocksUser "github.com/lechitz/AionApi/tests/mocks/user"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gorm.io/gorm"
)

type UserServiceTestSuite struct {
	Ctrl        *gomock.Controller
	Logger      *zap.SugaredLogger
	UserRepo    *mocksUser.MockRepository
	PasswordSvc *mocksSecurity.MockHasher
	TokenSvc    *mocksToken.MockStore
	UserSvc     *user.UserService
	Ctx         domain.ContextControl
}

func SetupUserServiceTest(t *testing.T) *UserServiceTestSuite {
	ctrl := gomock.NewController(t)

	logger := zaptest.NewLogger(t).Sugar()
	userRepo := mocksUser.NewMockRepository(ctrl)
	passwordSvc := mocksSecurity.NewMockHasher(ctrl)
	tokenSvc := mocksToken.NewMockStore(ctrl)

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
