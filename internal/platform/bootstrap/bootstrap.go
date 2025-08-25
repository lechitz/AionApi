// Package bootstrap composes secondary adapters and use cases (composition root).
package bootstrap

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache"
	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/tokenstore"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db"
	categoryRepo "github.com/lechitz/AionApi/internal/adapters/secondary/db/category/repository"
	userRepo "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security/hasher"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security/provider"

	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/core/usecase/auth"
	"github.com/lechitz/AionApi/internal/core/usecase/category"
	"github.com/lechitz/AionApi/internal/core/usecase/token"
	"github.com/lechitz/AionApi/internal/core/usecase/user"

	"github.com/lechitz/AionApi/internal/platform/bootstrap/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// AppDependencies exposes input ports that primary adapters consume.
type AppDependencies struct {
	TokenService    input.TokenService
	AuthService     input.AuthService
	UserService     input.UserService
	CategoryService input.CategoryService

	Logger output.ContextLogger
}

// InitializeDependencies wires infrastructure adapters and core use cases, returning input ports and a cleanup function.
func InitializeDependencies(appCtx context.Context, cfg *config.Config, logger output.ContextLogger) (*AppDependencies, func(context.Context), error) {
	// Infrastructure: cache
	cacheClient, err := cache.NewConnection(appCtx, cfg.Cache, logger)
	if err != nil {
		logger.Errorw(constants.ErrConnectToCache, commonkeys.Error, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgCacheConnected, commonkeys.CacheAddr, cfg.Cache.Addr)

	// Infrastructure: database
	dbConn, err := db.NewConnection(appCtx, cfg.DB, logger)
	if err != nil {
		logger.Errorw(constants.ErrConnectToDatabase, commonkeys.Error, err)
		return nil, nil, err
	}
	logger.Infow(constants.MsgPostgresConnected)

	// Secondary adapters (driven by output ports)
	passwordHasher := hasher.New()
	tokenProvider := provider.New(cfg.Secret.Key)
	tokenStore := tokenstore.New(cacheClient, logger)
	userRepository := userRepo.NewUser(dbConn, logger)
	categoryRepository := categoryRepo.NewCategory(dbConn, logger)

	// Core use cases (depend only on ports)
	tokenService := token.NewService(tokenStore, tokenProvider, logger)
	userService := user.NewService(userRepository, tokenStore, tokenProvider, passwordHasher, logger)
	authService := auth.NewService(userRepository, tokenStore, passwordHasher, tokenProvider, logger)
	categoryService := category.NewService(categoryRepository, logger)

	// Resource cleanup
	cleanup := func(ctx context.Context) {
		done := make(chan struct{})
		go func() {
			db.Close(dbConn, logger)
			if err := cacheClient.Close(); err != nil {
				logger.Errorw(constants.ErrCloseCacheConnection, commonkeys.Error, err)
			}
			close(done)
		}()
		select {
		case <-ctx.Done():
			logger.Warnw(constants.MsgCleanupAborted, commonkeys.Error, ctx.Err())
		case <-done:
			logger.Infow(constants.MsgCleanupCompletedSuccessfully)
		}
	}

	return &AppDependencies{
		TokenService:    tokenService,
		AuthService:     authService,
		UserService:     userService,
		CategoryService: categoryService,
		Logger:          logger,
	}, cleanup, nil
}
