// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
)

// ComposeRouter initializes and configures an HTTP router
// based on application dependencies and context path.
// Parameters:
//   - deps: application-wide dependencies
//   - contextPath: base prefix for all HTTP routes
//
// Returns:
//   - output.Router: the configured router
//   - error: any error encountered during setup
//
// ComposeRouter initializes and configures an HTTP router.
// The secretKey parameter is used to configure authentication middleware.
func ComposeRouter(deps *bootstrap.AppDependencies, contextPath, secretKey string) (output.Router, error) {
	httpRouter, err := NewHTTPRouter(deps.Logger, deps.TokenRepository, contextPath, secretKey)
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
		secretKey,
	)
}
