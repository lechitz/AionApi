// Package fxapp wires the application using Uber Fx modules.
package fxapp

import (
	"github.com/lechitz/AionApi/internal/adapter/secondary/hasher"
	"github.com/lechitz/AionApi/internal/adapter/secondary/token"
	authCache "github.com/lechitz/AionApi/internal/auth/adapter/secondary/cache"
	auth "github.com/lechitz/AionApi/internal/auth/core/usecase"
	categoryRepo "github.com/lechitz/AionApi/internal/category/adapter/secondary/db/repository"
	category "github.com/lechitz/AionApi/internal/category/core/usecase"
	chatClient "github.com/lechitz/AionApi/internal/chat/adapter/secondary/http"
	chat "github.com/lechitz/AionApi/internal/chat/core/usecase"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/cache"
	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/ports/output/httpclient"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	recordRepo "github.com/lechitz/AionApi/internal/record/adapter/secondary/db/repository"
	record "github.com/lechitz/AionApi/internal/record/core/usecase"
	tagRepo "github.com/lechitz/AionApi/internal/tag/adapter/secondary/db/repository"
	tag "github.com/lechitz/AionApi/internal/tag/core/usecase"
	userRepo "github.com/lechitz/AionApi/internal/user/adapter/secondary/db/repository"
	user "github.com/lechitz/AionApi/internal/user/core/usecase"
	"go.uber.org/fx"
)

// ApplicationModule wires the application layer (use cases, repositories, adapters) and exposes AppDependencies for HTTP composition.
//
//nolint:gochecknoglobals // Fx modules are intended as package-level options.
var ApplicationModule = fx.Options(
	fx.Provide(
		ProvideAppDependencies,
	),
)

// AppDependencies reuses the bootstrap contract for API composition.
type AppDependencies = bootstrap.AppDependencies

// ProvideAppDependencies composes repositories and use cases using shared infrastructure.
// Receives db.DB interface (not *gorm.DB) from InfraModule, following Dependency Inversion Principle.
func ProvideAppDependencies(cfg *config.Config, database db.DB, cacheClient cache.Cache, httpClient httpclient.HTTPClient, log logger.ContextLogger) *AppDependencies {
	passwordHasher := hasher.New()
	tokenProvider := token.NewProvider(cfg.Secret.Key)
	tokenStore := authCache.NewStore(cacheClient, log)

	userRepository := userRepo.New(database, log)
	categoryRepository := categoryRepo.New(database, log)
	tagRepository := tagRepo.New(database, log)
	recordRepository := recordRepo.New(database, log)

	chatHTTPClient := chatClient.NewClient(httpClient, cfg.AionChat.BaseURL, log)

	authService := auth.NewService(userRepository, tokenStore, tokenProvider, passwordHasher, log)
	userService := user.NewService(userRepository, tokenStore, tokenProvider, passwordHasher, log)
	categoryService := category.NewService(categoryRepository, log)
	tagService := tag.NewService(tagRepository, log)
	recordService := record.NewService(recordRepository, tagRepository, log)
	chatService := chat.NewService(chatHTTPClient, log)

	return &AppDependencies{
		AuthService:     authService,
		UserService:     userService,
		CategoryService: categoryService,
		TagService:      tagService,
		RecordService:   recordService,
		ChatService:     chatService,
		Logger:          log,
	}
}
