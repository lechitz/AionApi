// Package httpserver provides functionality for configuring and managing HTTP routers.
package httpserver

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"strings"

	"github.com/lechitz/AionApi/internal/adapters/primary/http/httpserver/constants"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/recovery"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/router/chi"

	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

// RouteComposer is a structure for configuring routes, middlewares, and logging in the HTTP router.
type RouteComposer struct {
	Router         portRouter.Router
	logger         logger.Logger
	authMiddleware *auth.MiddlewareAuth
	BasePath       string
}

// NewHTTPRouter creates and configures a new HTTP router with middleware and authentication.
// Parameters:
//   - logger: logger instance
//   - tokenRepository: token repository interface
//   - contextPath: base route path
//   - secretKey: JWT secret key (newly added dependency)
//
// Returns:
//   - *RouteComposer: configured route composer instance
//   - error: in case of any setup failure
func NewHTTPRouter(logger logger.Logger, tokenRepository output.TokenRepositoryPort, contextPath string, secretKey string) (*RouteComposer, error) {
	normalizedPath, err := normalizeContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()
	router.Use(recovery.RecoverMiddleware(logger))

	authMiddleware := auth.NewAuthMiddleware(tokenRepository, logger, secretKey)

	return &RouteComposer{
		BasePath:       normalizedPath,
		Router:         router,
		logger:         logger,
		authMiddleware: authMiddleware,
	}, nil
}

// GetRouter retrieves the current router instance used for managing HTTP routes.
func (r *RouteComposer) GetRouter() portRouter.Router {
	return r.Router
}

// normalizeContextPath ensures the given context path starts with '/' and is valid.
// Returns the normalized path or an error.
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
