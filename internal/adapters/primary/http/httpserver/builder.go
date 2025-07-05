// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/contextbuilder"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

// RouterBuilder is a struct for building and configuring HTTP routers
// with middleware and route handlers.
type RouterBuilder struct {
	Router      portRouter.Router
	ContextPath string
}

// BuildRouterRoutes sets up API routes, integrates middlewares,
// and returns the configured router with any error encountered.
// Parameters:
//   - logger: logging adapter
//   - userService: handles user-related use cases
//   - authService: handles authentication use cases
//   - tokenRepo: token cache storage interface
//   - contextPath: base path for route nesting
//   - adapter: router implementation
//
// Returns:
//   - portRouter.Router: configured router with routes and middleware
//   - error: any error encountered during setup
func BuildRouterRoutes(
	logger logger.Logger,
	userService http.UserService,
	authService input.AuthService,
	tokenRepo cache.TokenRepositoryPort,
	contextPath string,
	adapter portRouter.Router,
	secretKey string,
) (portRouter.Router, error) {
	adapter.Use(contextbuilder.InjectRequestIDMiddleware)

	genericHandler := handlers.NewGeneric(logger)
	userHandler := handlers.NewUser(userService, logger)
	authHandler := handlers.NewAuth(authService, logger)

	authMiddleware := auth.NewAuthMiddleware(tokenRepo, logger, secretKey)

	r := &RouteComposer{
		BasePath:       contextPath,
		Router:         adapter,
		logger:         logger,
		authMiddleware: authMiddleware,
	}

	adapter.Route(contextPath, func(rt portRouter.Router) {
		rt.Group(r.AddHealthCheckRoutes(genericHandler))
		rt.Group(r.AddUserRoutes(userHandler))
		rt.Group(r.AddAuthRoutes(authHandler))
	})

	return adapter, nil
}
