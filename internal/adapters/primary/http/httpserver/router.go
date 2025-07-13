// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lechitz/AionApi/internal/core/ports/input"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/platform/config"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/authmiddleware"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recoverymiddleware"
)

// RouteComposer is a structure for configuring routes, middlewares, and logging in the HTTP router.
type RouteComposer struct {
	Router         input.HTTPRouter
	logger         output.ContextLogger
	authMiddleware *authmiddleware.MiddlewareAuth
	BasePath       string
}

// NewHTTPRouter creates and configures a new HTTP router with middleware and authentication.
func NewHTTPRouter(
	tokenRepository output.TokenStore,
	cfg *config.Config,
	tokenClaimsExtractor output.TokenClaimsExtractor,
	logger output.ContextLogger,
	genericHandler *generic.Handler,
) (*RouteComposer, error) {
	normalizedPath, err := normalizeContextPath(cfg.ServerHTTP.Context)
	if err != nil {
		return nil, err
	}

	router := NewRouter()
	router.Use(recoverymiddleware.New(genericHandler))
	router.Use(injectRequestIDMiddleware)

	authMiddleware := authmiddleware.New(
		tokenRepository,
		logger,
		tokenClaimsExtractor,
	)

	return &RouteComposer{
		BasePath:       normalizedPath,
		Router:         router,
		logger:         logger,
		authMiddleware: authMiddleware,
	}, nil
}

// GetRouter retrieves the current router instance used for managing HTTP routes.
func (r *RouteComposer) GetRouter() input.HTTPRouter {
	return r.Router
}

// normalizeContextPath ensures the given context path starts with '/' and is valid.
func normalizeContextPath(raw string) (string, error) {
	if raw == "" {
		return "", errors.New(constants.ErrContextPathEmpty)
	}

	if strings.Contains(raw[1:], "/") {
		return "", errors.New(constants.ErrContextPathSlash)
	}

	if raw[0] != '/' {
		raw = "/" + raw
	}

	if len(raw) <= 1 {
		return "", fmt.Errorf(constants.InvalidContextPath, raw)
	}

	return raw, nil
}
