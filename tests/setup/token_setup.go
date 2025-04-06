package setup

import (
	"go.uber.org/zap"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
	"go.uber.org/zap/zaptest"
)

type TokenServiceTestSuite struct {
	Ctrl         *gomock.Controller
	TokenStore   *mockToken.MockTokenRepositoryPort
	TokenService token.TokenUsecase
	LoggerSugar  *zap.SugaredLogger
	Ctx          domain.ContextControl
}

func SetupTokenServiceTest(t *testing.T, secretKey string) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)
	logger := zaptest.NewLogger(t).Sugar()

	mockTokenStore := mockToken.NewMockTokenRepositoryPort(ctrl)

	tokenService := token.NewTokenService(mockTokenStore, logger, domain.TokenConfig{
		SecretKey: secretKey,
	})

	return &TokenServiceTestSuite{
		Ctrl:         ctrl,
		TokenStore:   mockTokenStore,
		TokenService: tokenService,
		LoggerSugar:  logger,
		Ctx:          domain.ContextControl{},
	}
}
