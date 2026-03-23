// Package main AionAPI
//
// @title AionAPI — REST API Documentation
// @version 0.0.1
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
	"path"

	swagger "github.com/lechitz/AionApi/contracts/openapi"
	"github.com/lechitz/AionApi/internal/platform/config"
)

const (
	// SwaggerTitle is the default human-readable title shown in the generated API documentation (Swagger/OpenAPI).
	SwaggerTitle = "AionAPI — REST API"
)

// configureSwagger applies runtime metadata from the loaded application configuration to the generated Swagger/OpenAPI spec.
func configureSwagger(cfg *config.Config) {
	swagger.SwaggerInfo.BasePath = path.Clean(cfg.ServerHTTP.Context + cfg.ServerHTTP.APIRoot)
	swagger.SwaggerInfo.Title = SwaggerTitle
	swagger.SwaggerInfo.Version = cfg.General.Version
}
