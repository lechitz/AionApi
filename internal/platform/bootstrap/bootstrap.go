// Package bootstrap provides a set of utilities for initializing application dependencies and managing application lifecycle.
package bootstrap

import (
	adapterCache "github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	infraCache "github.com/lechitz/AionApi/internal/adapters/secondary/cache/tools/redis"
	infraDB "github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/repository"
	adapterSecurity "github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/core/ports/input/graphql"
	"github.com/lechitz/AionApi/internal/core/ports/output/db"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/core/usecase/category"

	portsHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
	portsToken "github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/internal/platform/bootstrap/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// AppDependencies encapsulates all the core dependencies required for the application,
// including services, repositories, logging utilities and the loaded configuration.
type AppDependencies struct {
	Logger             logger.Logger                  // pointer
	TokenService       token.Usecase                  // pointer
	TokenRepository    portsToken.TokenRepositoryPort // pointer
	UserService        portsHttp.UserService          // pointer
	AuthService        portsHttp.AuthService          // pointer
	CategoryService    graphql.CategoryService        // pointer
	CategoryRepository db.CategoryStore               // pointer
	Config             config.Config                  // struct (não pointer)
}

// InitializeDependencies initializes and returns all core application dependencies,
// including repositories, services, and a cleanup function.
func InitializeDependencies(
	cfg config.Config,
	logger logger.Logger,
) (*AppDependencies, func(), error) {
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
	tokenService := token.NewTokenService(
		tokenRepository,
		logger,
		entity.TokenConfig{SecretKey: cfg.Secret.Key},
	)

	userRepository := repository.NewUserRepository(dbConn, logger)
	userService := user.NewUserService(userRepository, tokenService, passwordHasher, logger)

	categoryRepository := repository.NewCategoryRepository(dbConn, logger)
	categoryService := category.NewCategoryService(categoryRepository, logger)

	authService := auth.NewAuthService(
		userRepository,
		tokenService,
		passwordHasher,
		logger,
		cfg.Secret.Key,
	)

	cleanup := func() {
		infraDB.Close(dbConn, logger)

		if err := cacheConn.Close(); err != nil {
			logger.Errorf("%s: %v", constants.ErrCloseCacheConnection, err)
		}
	}

	return &AppDependencies{
		Config:             cfg,
		TokenRepository:    tokenRepository,
		TokenService:       tokenService,
		AuthService:        authService,
		UserService:        userService,
		CategoryRepository: categoryRepository,
		CategoryService:    categoryService,
		Logger:             logger,
	}, cleanup, nil
}
