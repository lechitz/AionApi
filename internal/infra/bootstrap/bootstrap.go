package bootstrap

import (
	adapterCache "github.com/lechitz/AionApi/adapters/secondary/cache"
	infraCache "github.com/lechitz/AionApi/adapters/secondary/cache/tools/redis"
	infraDB "github.com/lechitz/AionApi/adapters/secondary/db/postgres"
	adapterDB "github.com/lechitz/AionApi/adapters/secondary/db/repository"
	adapterSecurity "github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"

	portsHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	portsToken "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"

	"github.com/lechitz/AionApi/internal/infra/config"
)

type AppDependencies struct {
	TokenRepository portsToken.TokenRepositoryPort
	TokenService    token.TokenUsecase
	AuthService     portsHttp.AuthService
	UserService     portsHttp.UserService
	Logger          logger.Logger
}

const ErrorInitializingDependencies = "error closing cache connection"

func InitializeDependencies(cfg config.Config, logger logger.Logger) (*AppDependencies, func(), error) {

	cacheConn := infraCache.NewCacheConnection(cfg.Cache, logger)
	tokenRepository := adapterCache.NewTokenRepository(cacheConn, logger)

	dbConn := infraDB.NewDatabaseConnection(cfg.DB, logger)
	userRepository := adapterDB.NewUserRepository(dbConn, logger)

	tokenService := token.NewTokenService(tokenRepository, logger, domain.TokenConfig{
		SecretKey: cfg.SecretKey,
	})

	passwordHasher := adapterSecurity.NewBcryptPasswordAdapter()

	authService := auth.NewAuthService(userRepository, tokenService, passwordHasher, logger, cfg.SecretKey)
	userService := user.NewUserService(userRepository, tokenService, passwordHasher, logger)

	cleanup := func() {
		infraDB.Close(dbConn, logger)
		if err := cacheConn.Close(); err != nil {
			logger.Errorf("%s: %v", ErrorInitializingDependencies, err)
		}
	}

	return &AppDependencies{
		TokenRepository: tokenRepository,
		TokenService:    tokenService,
		AuthService:     authService,
		UserService:     userService,
		Logger:          logger,
	}, cleanup, nil
}
