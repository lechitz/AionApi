package setup

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/tests/mocks"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type TokenServiceTestSuite struct {
	Ctrl            *gomock.Controller
	TokenRepository *mocks.Moc
	TokenService    *token.TokenService
	LoggerSugar     *zap.SugaredLogger
	Ctx             domain.ContextControl
}

func SetupTokenServiceTest(t *testing.T) *TokenServiceTestSuite {
	ctrl := gomock.NewController(t)
	logger := zaptest.NewLogger(t).Sugar()

	client := mocks.NewMockCacheClient(ctrl)
	tokenRepo := cache.NewTokenRepository(client, logger)

	tokenService := token.NewTokenService(tokenRepo, logger, "XYZ1234567890")
	ctx := domain.ContextControl{}

	return &TokenServiceTestSuite{
		Ctrl:            ctrl,
		TokenRepository: tokenRepo,
		TokenService:    tokenService,
		LoggerSugar:     logger,
		Ctx:             ctx,
	}
}

var (
	TestPerfectToken = domain.TokenDomain{
		UserID: 1,
		Token:  "sampletoken",
	}
)
