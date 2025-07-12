// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// ComposeRouter initializes and configures an HTTP router based on application dependencies and context path.
func ComposeRouter(deps *bootstrap.AppDependencies, cfg *config.Config) (output.Router, error) {
	httpRouter, err := NewHTTPRouter(deps.TokenRepository, cfg.ServerHTTP.Context, deps.TokenClaimsExtractor, deps.Logger)
	if err != nil {
		return nil, err
	}

	return BuildRouterRoutes(
		deps.Logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenRepository,
		httpRouter.BasePath,
		httpRouter.GetRouter(),
		deps.TokenClaimsExtractor,
		cfg,
	)
}
