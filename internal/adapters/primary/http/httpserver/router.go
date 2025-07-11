// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/router/chi"
)

// RouteComposer is a structure for configuring routes, middlewares, and logging in the HTTP router.
type RouteComposer struct {
	Router         output.Router
	logger         output.Logger
	authMiddleware *auth.MiddlewareAuth
	BasePath       string
}

// NewHTTPRouter creates and configures a new HTTP router with middleware and authentication.
func NewHTTPRouter(
	logger output.Logger,
	tokenRepository output.TokenStore,
	contextPath string,
	tokenClaimsExtractor output.TokenClaimsExtractor,
) (*RouteComposer, error) {
	normalizedPath, err := normalizeContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(recovery.RecoverMiddleware(logger))

	authMiddleware := auth.NewAuthMiddleware(
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
func (r *RouteComposer) GetRouter() output.Router {
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
