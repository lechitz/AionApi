// Package main AionAPI
//
// @title AionAPI — REST API Documentation
// @version 0.1.0
// @description Public and internal REST API for AionAPI. Swagger (OpenAPI 3.x) is generated via swaggo.
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
// @x-graphQLPlayground /graphql/ // Cross-reference to GraphQL Playground (not covered by Swagger)
package main

import (
	"context"
	"os"
	"time"

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
		fx.Options(
			fx.Invoke(configureSwagger),
			fxapp.InfraModule,
			fxapp.ApplicationModule,
			fxapp.ServerModule,
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
