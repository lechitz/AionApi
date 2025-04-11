package bootstrap

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	portsToken "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	infraCache "github.com/lechitz/AionApi/internal/infrastructure/cache"
	infraDB "github.com/lechitz/AionApi/internal/infrastructure/db"

	adapterCache "github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	adapterDB "github.com/lechitz/AionApi/internal/adapters/secondary/db/repository"
	adapterSecurity "github.com/lechitz/AionApi/internal/infrastructure/security"

	portsHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"

	"github.com/lechitz/AionApi/internal/platform/config"

	"go.uber.org/zap"
)

type AppDependencies struct {
	TokenRepository portsToken.TokenRepositoryPort
	TokenService    token.TokenUsecase
	AuthService     portsHttp.AuthService
	UserService     portsHttp.UserService
}

var ErrorInitializingDependencies = "error closing cache connection"

func InitializeDependencies(logger *zap.SugaredLogger, cfg config.Config) (*AppDependencies, func(), error) {

	cacheConn := infraCache.NewCacheConnection(cfg.Cache, logger)
	tokenRepository := adapterCache.NewTokenRepository(cacheConn, logger)

	dbConn := infraDB.NewDatabaseConnection(cfg.DB, logger)
	userRepository := adapterDB.NewUserRepository(dbConn, logger)

	tokenService := token.NewTokenService(tokenRepository, logger, domain.TokenConfig{
		SecretKey: cfg.SecretKey,
	})

	authService := auth.NewAuthService(userRepository, tokenService, adapterSecurity.BcryptPasswordAdapter{}, logger, cfg.SecretKey)
	userService := user.NewUserService(userRepository, tokenService, adapterSecurity.BcryptPasswordAdapter{}, logger)

	cleanup := func() {
		infraDB.Close(dbConn, logger)
		if err := cacheConn.Close(); err != nil {
			logger.Error(ErrorInitializingDependencies, err)
		}
	}

	return &AppDependencies{
		TokenRepository: tokenRepository,
		TokenService:    tokenService,
		AuthService:     authService,
		UserService:     userService,
	}, cleanup, nil
}
