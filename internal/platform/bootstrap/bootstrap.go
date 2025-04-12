package bootstrap

import (
	adapterCache "github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	adapterDB "github.com/lechitz/AionApi/internal/adapters/secondary/db/repository"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"

	portsToken "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	infraCache "github.com/lechitz/AionApi/internal/infrastructure/cache"
	infraDB "github.com/lechitz/AionApi/internal/infrastructure/db"

	adapterSecurity "github.com/lechitz/AionApi/internal/infrastructure/security"

	portsHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"

	"github.com/lechitz/AionApi/internal/platform/config"
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
