// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
)

// ComposeRouter initializes and configures an HTTP router
// based on application dependencies and context path.
// Parameters:
//   - deps: application-wide dependencies
//   - contextPath: base prefix for all HTTP routes
//
// Returns:
//   - portRouter.Router: the configured router
//   - error: any error encountered during setup
func ComposeRouter(deps *bootstrap.AppDependencies, contextPath string) (portRouter.Router, error) {
	httpRouter, err := NewHTTPRouter(deps.Logger, deps.TokenRepository, contextPath, deps.Config.Secret.Key)
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
		deps.Config.Secret.Key,
	)
}
