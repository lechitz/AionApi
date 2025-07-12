package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// TokenServiceTestSuite is a test suite for testing TokenService with mocked dependencies and utilities for test cases.// TokenServiceTestSuite is a struct for managing the dependencies needed to test the TokenService implementation.
type TokenServiceTestSuite struct {
	Ctx          context.Context
	TokenService *token.Service
	Ctrl         *gomock.Controller
	Logger       *mocks.MockLogger
	TokenStore   *mocks.MockTokenStore
}

// TokenServiceTest initializes a test suite for Service with mocked dependencies and a given secret key.
func TokenServiceTest(t *testing.T, secretKey config.Secret) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockLog := mocks.NewMockLogger(ctrl)
	mockTokenStore := mocks.NewMockTokenStore(ctrl)

	ExpectLoggerDefaultBehavior(mockLog)

	tokenService := token.NewService(mockTokenStore, mockLog, secretKey)

	return &TokenServiceTestSuite{
		Ctrl:         ctrl,
		Logger:       mockLog,
		TokenStore:   mockTokenStore,
		TokenService: tokenService,
		Ctx:          t.Context(),
	}
}
