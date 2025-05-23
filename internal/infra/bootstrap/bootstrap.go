package bootstrap

import (
	adapterCache "github.com/lechitz/AionApi/adapters/secondary/cache"
	infraCache "github.com/lechitz/AionApi/adapters/secondary/cache/tools/redis"
	infraDB "github.com/lechitz/AionApi/adapters/secondary/db/postgres"
	adapterDB "github.com/lechitz/AionApi/adapters/secondary/db/repository"
	adapterSecurity "github.com/lechitz/AionApi/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/input/graphql"
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/usecase/category"
	"github.com/lechitz/AionApi/internal/infra/bootstrap/constants"

	portsHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	portsToken "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/internal/infra/config"
)

type AppDependencies struct {
	TokenRepository    portsToken.TokenRepositoryPort
	TokenService       token.TokenUsecase
	AuthService        portsHttp.AuthService
	UserService        portsHttp.UserService
	CategoryService    graphql.CategoryService
	CategoryRepository db.CategoryStore
	Logger             logger.Logger
}

func InitializeDependencies(cfg config.Config, logger logger.Logger) (*AppDependencies, func(), error) {

	cacheConn, err := infraCache.NewCacheConnection(cfg.Cache, logger)
	if err != nil {
		logger.Errorf(constants.ErrConnectToCache, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgCacheConnected, constants.FieldAddr, cfg.Cache.Addr)

	dbConn, err := infraDB.NewDatabaseConnection(cfg.DB, logger)
	if err != nil {
		logger.Errorf(constants.ErrConnectToDatabase, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgPostgresConnected)

	passwordHasher := adapterSecurity.NewBcryptPasswordAdapter()

	tokenRepository := adapterCache.NewTokenRepository(cacheConn, logger)
	tokenService := token.NewTokenService(tokenRepository, logger, domain.TokenConfig{SecretKey: cfg.Secret.Key})

	userRepository := adapterDB.NewUserRepository(dbConn, logger)
	userService := user.NewUserService(userRepository, tokenService, passwordHasher, logger)

	categoryRepository := adapterDB.NewCategoryRepository(dbConn, logger)
	categoryService := category.NewCategoryService(categoryRepository, logger)

	authService := auth.NewAuthService(userRepository, tokenService, passwordHasher, logger, cfg.Secret.Key)

	cleanup := func() {
		infraDB.Close(dbConn, logger)

		if err := cacheConn.Close(); err != nil {
			logger.Errorf("%s: %v", constants.ErrCloseCacheConnection, err)
		}
	}

	return &AppDependencies{
		TokenRepository:    tokenRepository,
		TokenService:       tokenService,
		AuthService:        authService,
		UserService:        userService,
		CategoryRepository: categoryRepository,
		CategoryService:    categoryService,
		Logger:             logger,
	}, cleanup, nil
}
