package setup

import (
	"context"
	"testing"

	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/mock/gomock"
)

// TokenServiceTestSuite groups mocked dependencies and the SUT (TokenService)
// to simplify Token-related unit tests.
type TokenServiceTestSuite struct {
	Ctx           context.Context
	TokenService  *token.Service
	Ctrl          *gomock.Controller
	Logger        *mocks.ContextLogger
	TokenStore    *mocks.TokenStore
	TokenProvider *mocks.TokenProvider
}

// TokenServiceTest initializes and returns a TokenServiceTestSuite with mocked output ports.
func TokenServiceTest(t *testing.T) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)

	logger := mocks.NewContextLogger(ctrl)
	tokenStore := mocks.NewTokenStore(ctrl)
	tokenProvider := mocks.NewTokenProvider(ctrl)

	ExpectLoggerDefaultBehavior(logger)

	svc := token.NewService(tokenStore, tokenProvider, logger)

	return &TokenServiceTestSuite{
		Ctrl:          ctrl,
		Logger:        logger,
		TokenStore:    tokenStore,
		TokenProvider: tokenProvider,
		TokenService:  svc,
		Ctx:           t.Context(),
	}
}
