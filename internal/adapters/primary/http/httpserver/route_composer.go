// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// ComposeRouter initializes and configures an HTTP router based on application dependencies and context path.
func ComposeRouter(deps *bootstrap.AppDependencies, cfg *config.Config) (input.HTTPRouter, error) {
	genericHandler := generic.New(deps.Logger, cfg.General)

	httpRouter, err := NewHTTPRouter(
		deps.TokenRepository,
		cfg,
		deps.TokenClaimsExtractor,
		deps.Logger,
		genericHandler,
	)
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
		genericHandler,
	)
}
