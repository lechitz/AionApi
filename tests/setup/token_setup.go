package setup

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	mockLogger "github.com/lechitz/AionApi/tests/mocks/logger"
	mockToken "github.com/lechitz/AionApi/tests/mocks/token"
)

type TokenServiceTestSuite struct {
	Ctrl         *gomock.Controller
	Logger       *mockLogger.MockLogger
	TokenStore   *mockToken.MockTokenRepositoryPort
	TokenService token.TokenUsecase
	Ctx          context.Context
}

func SetupTokenServiceTest(t *testing.T, secretKey string) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockLog := mockLogger.NewMockLogger(ctrl)
	mockTokenStore := mockToken.NewMockTokenRepositoryPort(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	tokenService := token.NewTokenService(mockTokenStore, mockLog, domain.TokenConfig{
		SecretKey: secretKey,
	})

	return &TokenServiceTestSuite{
		Ctrl:         ctrl,
		Logger:       mockLog,
		TokenStore:   mockTokenStore,
		TokenService: tokenService,
		Ctx:          context.Background(),
	}
}
