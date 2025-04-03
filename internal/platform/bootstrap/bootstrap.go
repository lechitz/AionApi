package bootstrap

import (
	"github.com/lechitz/AionApi/config"
	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	dbadapter "github.com/lechitz/AionApi/internal/adapters/secondary/db"
	httpports "github.com/lechitz/AionApi/internal/core/ports/input/http"
	tokenports "github.com/lechitz/AionApi/internal/core/ports/output/token"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	tokeninfra "github.com/lechitz/AionApi/internal/infrastructure/cache"
	"github.com/lechitz/AionApi/internal/infrastructure/db/postgres"
	securityadapter "github.com/lechitz/AionApi/internal/infrastructure/security"

	"go.uber.org/zap"
)

type AppDependencies struct {
	UserService  httpports.UserService
	AuthService  httpports.AuthService
	TokenService tokenports.Store
}

var ErrorInitializingDependencies = "error closing cache connection"

func InitializeDependencies(logger *zap.SugaredLogger, cfg config.Config) (*AppDependencies, func(), error) {

	redisClient := tokeninfra.NewRedisConnection(cfg.CacheConfig, logger)
	tokenStore := cache.NewTokenRepository(redisClient, logger)

	dbConn := postgres.NewDatabaseConnection(cfg.DBConfig, logger)
	userRepo := dbadapter.NewUserRepository(dbConn, logger)

	tokenService := token.NewTokenService(*tokenStore, logger, cfg.SecretKey)
	userService := user.NewUserService(userRepo, tokenService, securityadapter.BcryptPasswordAdapter{}, logger)
	authService := auth.NewAuthService(userRepo, tokenService, securityadapter.BcryptPasswordAdapter{}, logger, cfg.SecretKey)

	cleanup := func() {
		postgres.Close(dbConn, logger)
		if err := redisClient.Close(); err != nil {
			logger.Error(ErrorInitializingDependencies, err)
		}
	}

	return &AppDependencies{
		UserService:  userService,
		AuthService:  authService,
		TokenService: tokenService,
	}, cleanup, nil
}
