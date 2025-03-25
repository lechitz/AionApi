package bootstrap

import (
	tokenadapter "github.com/lechitz/AionApi/adapters/output/cache"
	dbadapter "github.com/lechitz/AionApi/adapters/output/db"
	securityadapter "github.com/lechitz/AionApi/adapters/output/security"

	"github.com/lechitz/AionApi/app/config"
	"github.com/lechitz/AionApi/core/service"
	"github.com/lechitz/AionApi/infra/cache"
	"github.com/lechitz/AionApi/infra/db"

	inputHttp "github.com/lechitz/AionApi/ports/input/http"
	portSecurity "github.com/lechitz/AionApi/ports/output/security"

	"go.uber.org/zap"
)

type AppDependencies struct {
	UserService  inputHttp.IUserService
	AuthService  inputHttp.IAuthService
	TokenService portSecurity.ITokenService

	Config config.Config
}

const ErrorInitializingDependencies = "error closing cache connection"

func InitializeDependencies(loggerSugar *zap.SugaredLogger, cfg config.Config) (*AppDependencies, func(), error) {

	cacheConn := cache.NewCacheConnection(cfg.CacheConfig, loggerSugar)
	tokenStore := tokenadapter.NewTokenStore(cacheConn, loggerSugar)

	databaseConn := db.NewDatabaseConnection(cfg.DBConfig, loggerSugar)
	userRepo := dbadapter.NewUserRepo(databaseConn, loggerSugar)

	var passwordHasher portSecurity.IPasswordService = securityadapter.BcryptPasswordAdapter{}
	var tokenService portSecurity.ITokenService = service.NewTokenService(tokenStore, loggerSugar, cfg.SecretKey)

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
