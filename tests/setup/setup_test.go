package setup

import (
	"github.com/lechitz/AionApi/internal/core/service/user"
	"go.uber.org/zap"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/zap/zaptest"
)

type TestDependencies struct {
	Controller       *gomock.Controller
	Logger           *zap.SugaredLogger
	MockUserRepo     *mocks.MockIUserRepository
	MockTokenService *mocks.MockITokenService
	MockPasswordSvc  *mocks.MockIPasswordService
	Service          *user.UserService
}

func NewUserServiceTestSetup(t *testing.T) *TestDependencies {
	ctrl := gomock.NewController(t)
	logger := zaptest.NewLogger(t).Sugar()

	mockUserRepo := mocks.NewMockIUserRepository(ctrl)
	mockTokenSvc := mocks.NewMockITokenService(ctrl)
	mockPasswordSvc := mocks.NewMockIPasswordService(ctrl)

	userService := user.NewUserService(
		mockUserRepo,
		mockTokenSvc,
		mockPasswordSvc,
		logger,
	)

	return &TestDependencies{
		Controller:       ctrl,
		Logger:           logger,
		MockUserRepo:     mockUserRepo,
		MockTokenService: mockTokenSvc,
		MockPasswordSvc:  mockPasswordSvc,
		Service:          userService,
	}
}
