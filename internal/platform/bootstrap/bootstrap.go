// Package bootstrap provides a set of utilities for initializing application dependencies and managing the application lifecycle.
package bootstrap

import (
	"context"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	adapterToken "github.com/lechitz/AionApi/internal/adapters/secondary/cache/token"
	adapterCache "github.com/lechitz/AionApi/internal/adapters/secondary/cache/tools/redis"
	adapterDB "github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/repository"
	adapterSecurity "github.com/lechitz/AionApi/internal/adapters/secondary/security"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/core/usecase/category"

	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"
	"github.com/lechitz/AionApi/internal/platform/bootstrap/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// AppDependencies bundles the services and adapters used throughout the application.
type AppDependencies struct {
	AuthService     input.AuthService
	UserService     input.UserService
	CategoryService input.CategoryService
	TokenService    input.TokenService

	TokenClaimsExtractor output.TokenClaimsExtractor
	TokenRepository      output.TokenStore
	Logger               output.Logger
}

// InitializeDependencies initializes and returns all core application dependencies.
func InitializeDependencies(appCtx context.Context, cfg *config.Config, logger output.Logger) (*AppDependencies, func(), error) {
	cacheClient, err := adapterCache.NewCacheConnection(appCtx, cfg.Cache, logger)
	if err != nil {
		logger.Errorf(constants.ErrConnectToCache, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgCacheConnected, commonkeys.CacheAddr, cfg.Cache.Addr)

	dbConn, err := adapterDB.NewDatabaseConnection(appCtx, cfg.DB, logger)
	if err != nil {
		logger.Errorf(constants.ErrConnectToDatabase, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgPostgresConnected)

	// Security Hasher
	passwordHasher := adapterSecurity.NewBcryptPasswordAdapter()

	// Token Extractor
	tokenClaimsExtractor := adapterSecurity.NewJWTClaimsExtractor(cfg.Secret.Key)

	// Token
	tokenRepository := adapterToken.NewTokenRepository(cacheClient, logger)
	tokenService := token.NewTokenService(tokenRepository, logger, cfg.Secret)

	// User
	userRepository := repository.NewUserRepository(dbConn, logger)
	userService := user.NewUserService(userRepository, tokenService, passwordHasher, logger)

	// Category
	categoryRepository := repository.NewCategoryRepository(dbConn, logger)
	categoryService := category.NewCategoryService(categoryRepository, logger)

	// Auth
	authService := auth.NewAuthService(userRepository, tokenService, passwordHasher, logger)

	cleanupResources := func() {
		adapterDB.Close(dbConn, logger)

		if err := cacheClient.Close(); err != nil {
			logger.Errorf("%s: %v", constants.ErrCloseCacheConnection, err)
		}
	}

	return &AppDependencies{
		AuthService:     authService,
		UserService:     userService,
		CategoryService: categoryService,
		TokenService:    tokenService,

		TokenClaimsExtractor: tokenClaimsExtractor,
		TokenRepository:      tokenRepository,
		Logger:               logger,
	}, cleanupResources, nil
}
