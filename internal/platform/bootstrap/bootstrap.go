// Package bootstrap composes secondary adapter and use cases (composition root).
package bootstrap

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapter/secondary/cache/redis"
	"github.com/lechitz/AionApi/internal/adapter/secondary/db/postgres"
	"github.com/lechitz/AionApi/internal/adapter/secondary/hasher"
	"github.com/lechitz/AionApi/internal/adapter/secondary/token"
	"github.com/lechitz/AionApi/internal/auth/adapter/secondary/cache"
	inputAuth "github.com/lechitz/AionApi/internal/auth/core/ports/input"
	auth "github.com/lechitz/AionApi/internal/auth/core/usecase"
	categoryRepo "github.com/lechitz/AionApi/internal/category/adapter/secondary/db/repository"
	inputCategory "github.com/lechitz/AionApi/internal/category/core/ports/input"
	category "github.com/lechitz/AionApi/internal/category/core/usecase"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	recordRepo "github.com/lechitz/AionApi/internal/record/adapter/secondary/db/repository"
	inputRecord "github.com/lechitz/AionApi/internal/record/core/ports/input"
	record "github.com/lechitz/AionApi/internal/record/core/usecase"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	tagRepo "github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/repository"
	inputTag "github.com/lechitz/AionApi/internal/tag/core/ports/input"
	tag "github.com/lechitz/AionApi/internal/tag/core/usecase"
	userRepo "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/repository"
	inputUser "github.com/lechitz/AionApi/internal/user/core/ports/input"
	user "github.com/lechitz/AionApi/internal/user/core/usecase"
)

// AppDependencies exposes input ports that the primary adapter consumes.
type AppDependencies struct {
	AuthService     inputAuth.AuthService
	UserService     inputUser.UserService
	CategoryService inputCategory.CategoryService
	TagService      inputTag.TagService
	RecordService   inputRecord.RecordService

	Logger logger.ContextLogger
}

// InitializeDependencies wires infrastructure adapter and core use cases, returning input ports and a cleanup function.
func InitializeDependencies(appCtx context.Context, cfg *config.Config, logger logger.ContextLogger) (*AppDependencies, func(context.Context), error) {
	// Infrastructure: cache
	cacheClient, err := redis.NewConnection(appCtx, cfg.Cache, logger)
	if err != nil {
		logger.Errorw(ErrConnectToCache, commonkeys.Error, err)
		return nil, nil, err
	}
	logger.Infow(MsgCacheConnected, commonkeys.CacheAddr, cfg.Cache.Addr)

	// Infrastructure: database
	dbConn, err := postgres.NewConnection(appCtx, cfg.DB, logger)
	if err != nil {
		logger.Errorw(ErrConnectToDatabase, commonkeys.Error, err)
		return nil, nil, err
	}
	logger.Infow(MsgPostgresConnected)

	// Secondary adapter (driven by output ports)
	passwordHasher := hasher.New()
	tokenProvider := token.NewProvider(cfg.Secret.Key)
	tokenStore := cache.NewStore(cacheClient, logger)
	userRepository := userRepo.New(dbConn, logger)
	categoryRepository := categoryRepo.New(dbConn, logger)
	tagRepository := tagRepo.New(dbConn, logger)
	recordRepository := recordRepo.New(dbConn, logger)

	// Core use cases (depend only on ports)
	authService := auth.NewService(userRepository, tokenStore, tokenProvider, passwordHasher, logger)
	userService := user.NewService(userRepository, tokenStore, tokenProvider, passwordHasher, logger)
	categoryService := category.NewService(categoryRepository, logger)
	tagService := tag.NewService(tagRepository, logger)
	recordService := record.NewService(recordRepository, tagRepository, logger)

	// Resource cleanup
	cleanup := func(ctx context.Context) {
		done := make(chan struct{})
		go func() {
			postgres.Close(dbConn, logger)
			if err := cacheClient.Close(); err != nil {
				logger.Errorw(ErrCloseCacheConnection, commonkeys.Error, err)
			}
			close(done)
		}()
		select {
		case <-ctx.Done():
			logger.Warnw(MsgCleanupAborted, commonkeys.Error, ctx.Err())
		case <-done:
			logger.Infow(MsgCleanupCompletedSuccessfully)
		}
	}

	return &AppDependencies{
		AuthService:     authService,
		UserService:     userService,
		CategoryService: categoryService,
		TagService:      tagService,
		RecordService:   recordService,
		Logger:          logger,
	}, cleanup, nil
}
