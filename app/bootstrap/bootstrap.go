package bootstrap

import (
	"github.com/lechitz/AionApi/adapters/output/cache"
	dbadapter "github.com/lechitz/AionApi/adapters/output/db"
	portsCache "github.com/lechitz/AionApi/ports/output/cache"
	portsDB "github.com/lechitz/AionApi/ports/output/db"

	"github.com/lechitz/AionApi/app/config"
	"github.com/lechitz/AionApi/core/service"
	infraCache "github.com/lechitz/AionApi/infra/cache"
	infraDB "github.com/lechitz/AionApi/infra/db"
	"go.uber.org/zap"
)

type AppDependencies struct {
	TokenService *service.TokenService
	UserService  *service.UserService
	AuthService  *service.AuthService

	TokenRepository portsCache.ITokenRepository
	UserRepository  portsDB.IUserRepository

	Config config.Config
}

const ErrorInitializingDependencies = "error closing cache connection: "

func InitializeDependencies(loggerSugar *zap.SugaredLogger, cfg config.Config) (*AppDependencies, func(), error) {
	cacheConn := infraCache.NewCacheConnection(cfg.CacheConfig, loggerSugar)
	tokenRepo := cache.NewCacheRepo(cacheConn, loggerSugar)

	databaseConn := infraDB.NewDatabaseConnection(cfg.DBConfig, loggerSugar)
	userRepo := dbadapter.NewUserRepo(databaseConn, loggerSugar)

	tokenService := service.NewTokenService(userRepo, tokenRepo, loggerSugar, cfg.SecretKey)
	authService := service.NewAuthService(userRepo, tokenRepo, tokenService, loggerSugar, cfg.SecretKey)
	userService := service.NewUserService(userRepo, tokenService, authService, loggerSugar)

	cleanup := func() {
		infraDB.Close(databaseConn, loggerSugar)
		if err := cacheConn.Close(); err != nil {
			loggerSugar.Error(ErrorInitializingDependencies, err)
		}
	}

	return &AppDependencies{
		TokenService:    tokenService,
		TokenRepository: tokenRepo,
		UserService:     userService,
		AuthService:     authService,
		UserRepository:  userRepo,
		Config:          cfg,
	}, cleanup, nil
}
