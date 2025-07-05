package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/domain/entity"

	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// TokenServiceTestSuite is a test suite for testing TokenService with mocked dependencies and utilities for test cases.// TokenServiceTestSuite is a struct for managing the dependencies needed to test the TokenService implementation.
type TokenServiceTestSuite struct {
	Ctrl         *gomock.Controller
	Logger       *mocks.MockLogger
	TokenStore   *mocks.MockTokenRepositoryPort
	TokenService token.Usecase
	Ctx          context.Context
}

// TokenServiceTest initializes a test suite for Service with mocked dependencies and a given secret key.
func TokenServiceTest(t *testing.T, secretKey string) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockLog := mocks.NewMockLogger(ctrl)
	mockTokenStore := mocks.NewMockTokenRepositoryPort(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	tokenService := token.NewTokenService(mockTokenStore, mockLog, entity.TokenConfig{
		SecretKey: secretKey,
	})

	return &TokenServiceTestSuite{
		Ctrl:         ctrl,
		Logger:       mockLog,
		TokenStore:   mockTokenStore,
		TokenService: tokenService,
		Ctx:          t.Context(),
	}
}
