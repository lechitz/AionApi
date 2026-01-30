package main

import (
	"path"

	swagger "github.com/lechitz/AionApi/docs/swagger"
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
