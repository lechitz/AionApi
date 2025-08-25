// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/handler"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// ComposeRouter builds the router and wires routes/middlewares.
func ComposeRouter(cfg *config.Config, tokenSvc input.TokenService, userService input.UserService, authService input.AuthService, logger output.ContextLogger) (http.Handler, error) {
	genericHandler := handler.New(logger, cfg.General)

	httpRouter, err := NewHTTPRouter(cfg, tokenSvc, genericHandler, logger)
	if err != nil {
		return nil, err
	}

	if err := httpRouter.WireRoutes(cfg, userService, authService, genericHandler, logger); err != nil {
		return nil, err
	}

	// Expose the concrete router as http.Handler
	if h, ok := httpRouter.Router.(http.Handler); ok {
		return h, nil
	}
	return nil, nil
}
