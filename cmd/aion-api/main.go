// Package main AionAPI
//
// @title           AionAPI — REST API Documentation
// @version         0.1.0
// @description     Public and internal REST API for AionAPI. Swagger (OpenAPI 3.x) is generated via swaggo.
// @termsOfService  https://github.com/lechitz/AionApi
//
// @contact.name   Lechitz
// @contact.url    https://github.com/lechitz
// @contact.email  felipe.lechitz@gmail.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @schemes http https
// @BasePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Use "Bearer {token}" (JWT).
//
// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name auth_token
// @description Session-based authentication cookie (e.g., Redis-backed).
//
// @x-graphQLPlayground /graphql/  // Cross-reference to GraphQL Playground (not covered by Swagger)
package main

import (
	"context"
	"os"
	"path"
	"time"

	swagger "github.com/lechitz/AionApi/docs/swagger"
	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/fxapp"
	"go.uber.org/fx"
)

// the main is the entry point for the application.
func main() {
	os.Exit(run())
}

// run is the main application logic.
func run() int {
	app := fx.New(
		fx.Provide(
			fxapp.ProvideLogger,
			fxapp.ProvideKeyGenerator,
			fxapp.ProvideConfig,
			fxapp.ProvideDependencies,
		),
		fx.Invoke(
			configureSwagger,
			fxapp.InitObservability,
			fxapp.RunServers,
		),
	)
	if err := app.Start(context.Background()); err != nil {
		return 1
	}
	<-app.Done()

	stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Stop(stopCtx); err != nil {
		return 1
	}
	return 0
}

// configureSwagger applies runtime metadata from the loaded application configuration to the generated Swagger/OpenAPI spec.
func configureSwagger(cfg *config.Config) {
	swagger.SwaggerInfo.BasePath = path.Clean(cfg.ServerHTTP.Context + cfg.ServerHTTP.APIRoot)
	swagger.SwaggerInfo.Title = SwaggerTitle
	swagger.SwaggerInfo.Version = cfg.General.Version
}
