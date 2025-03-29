package bootstrap

import (
	tokenadapter "github.com/lechitz/AionApi/adapters/output/cache"
	dbadapter "github.com/lechitz/AionApi/adapters/output/db"
	securityadapter "github.com/lechitz/AionApi/adapters/output/security"
	"github.com/lechitz/AionApi/core/ports/input/http"
	"github.com/lechitz/AionApi/core/ports/output/security"
	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/core/service"
	"github.com/lechitz/AionApi/infra/cache"
	"github.com/lechitz/AionApi/infra/db"

	"go.uber.org/zap"
)

type AppDependencies struct {
	UserService  http.IUserService
	AuthService  http.IAuthService
	TokenService security.ITokenService

	Config config.Config
}

const ErrorInitializingDependencies = "error closing cache connection"

func InitializeDependencies(loggerSugar *zap.SugaredLogger, cfg config.Config) (*AppDependencies, func(), error) {

	cacheConn := cache.NewCacheConnection(cfg.CacheConfig, loggerSugar)
	tokenStore := tokenadapter.NewTokenStore(cacheConn, loggerSugar)

	databaseConn := db.NewDatabaseConnection(cfg.DBConfig, loggerSugar)
	userRepo := dbadapter.NewUserRepo(databaseConn, loggerSugar)

	var passwordHasher security.IPasswordService = securityadapter.BcryptPasswordAdapter{}
	var tokenService security.ITokenService = service.NewTokenService(tokenStore, loggerSugar, cfg.SecretKey)

	userService := service.NewUserService(userRepo, tokenService, passwordHasher, loggerSugar)
	authService := service.NewAuthService(userRepo, tokenService, passwordHasher, loggerSugar, cfg.SecretKey)

	cleanup := func() {
		db.Close(databaseConn, loggerSugar)
		if err := cacheConn.Close(); err != nil {
			loggerSugar.Error(ErrorInitializingDependencies, err)
		}
	}

	return &AppDependencies{
		TokenService: tokenService,
		UserService:  userService,
		AuthService:  authService,
		Config:       cfg,
	}, cleanup, nil
}
