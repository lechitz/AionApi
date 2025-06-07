package httpserver

import (
	"fmt"
	"strings"

	"github.com/lechitz/AionApi/adapters/primary/http/httpserver/constants"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/recovery"
	"github.com/lechitz/AionApi/adapters/primary/http/router/chi"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

type RouteComposer struct {
	BasePath       string
	Router         portRouter.Router
	logger         logger.Logger
	authMiddleware *auth.MiddlewareAuth
}

func NewHttpRouter(logger logger.Logger, tokenRepository cache.TokenRepositoryPort, contextPath string) (*RouteComposer, error) {
	normalizedPath, err := normalizeContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	router := chi.NewRouter()

	router.Use(recovery.RecoverMiddleware(logger))

	authMiddleware := auth.NewAuthMiddleware(tokenRepository, logger)

	return &RouteComposer{
		BasePath:       normalizedPath,
		Router:         router,
		logger:         logger,
		authMiddleware: authMiddleware,
	}, nil
}

func (r *RouteComposer) GetRouter() portRouter.Router {
	return r.Router
}

func normalizeContextPath(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf(constants.ErrContextPathEmpty)
	}

	if strings.Contains(raw[1:], "/") {
		return "", fmt.Errorf(constants.ErrContextPathSlash)
	}

	if raw[0] != '/' {
		raw = "/" + raw
	}

	if len(raw) <= 1 {
		return "", fmt.Errorf(constants.InvalidContextPath, raw)
	}

	return raw, nil
}
