// Package httpserver wires route groups into the router.
package httpserver

import (
	handlerAuth "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/auth/handler"
	handlerGeneric "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/handler"
	handlerUser "github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/user/handler"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// WireRoutes registers all routes into this RouteComposer.
func (r *RouteComposer) WireRoutes(cfg *config.Config, userService input.UserService, authService input.AuthService, genericHandler *handlerGeneric.Handler, logger output.ContextLogger) error {
	userHandler := handlerUser.New(userService, cfg, logger)
	authHandler := handlerAuth.New(authService, cfg, logger)

	r.Router.SetNotFoundHandler(genericHandler.NotFoundHandler)
	r.Router.SetMethodNotAllowedHandler(genericHandler.MethodNotAllowedHandler)
	r.Router.SetErrorHandler(genericHandler.ErrorHandler)

	r.Router.Route(r.BasePath, func(rt HTTPRouter) {
		rt.Group(r.AddHealthCheckRoutes(genericHandler))
		rt.Group(r.AddUserRoutes(userHandler))
		rt.Group(r.AddAuthRoutes(authHandler))
	})
	return nil
}
