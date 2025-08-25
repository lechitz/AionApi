// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/controllers/generic/handler"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recoverymiddleware"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/requestidmiddleware"
	"github.com/lechitz/AionApi/internal/core/ports/input"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/platform/config"
)

// HTTPRouter is a minimal interface to compose routes and middlewares.
type HTTPRouter interface {
	Use(middleware func(http.Handler) http.Handler)
	Route(pattern string, fn func(r HTTPRouter))
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
	Mount(pattern string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Group(fn func(r HTTPRouter))
	SetNotFoundHandler(handler http.HandlerFunc)
	SetMethodNotAllowedHandler(handler http.HandlerFunc)
	SetErrorHandler(handler func(http.ResponseWriter, *http.Request, error))
	GroupWithMiddleware(middleware func(http.Handler) http.Handler, fn func(r HTTPRouter))
}

// RouteComposer wires router, base path and auth middleware.
type RouteComposer struct {
	Router         HTTPRouter
	logger         output.ContextLogger
	authMiddleware *authmiddleware.MiddlewareAuth
	BasePath       string
}

// NewHTTPRouter sets up recovery/request-id middlewares and auth middleware.
// NOTE: Auth middleware only extracts the raw token and delegates validation to TokenService.
func NewHTTPRouter(cfg *config.Config, tokenSvc input.TokenService, genericHandler *handler.Handler, logger output.ContextLogger) (*RouteComposer, error) {
	router := NewRouter()

	// Generic middlewares
	router.Use(recoverymiddleware.New(genericHandler))
	router.Use(requestidmiddleware.New())

	// Auth middleware (no JWT parsing here)
	amw := authmiddleware.New(tokenSvc, logger)

	return &RouteComposer{
		BasePath:       cfg.ServerHTTP.Context,
		Router:         router,
		logger:         logger,
		authMiddleware: amw,
	}, nil
}

// GetRouter exposes the underlying router (useful for tests/customization).
func (r *RouteComposer) GetRouter() HTTPRouter {
	return r.Router
}
