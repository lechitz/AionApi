package server

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/server/constants"
	routerAdapter "github.com/lechitz/AionApi/internal/adapters/primary/routeradapter"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
	"strings"
)

type HttpRouter struct {
	ContextPath    string
	Router         portRouter.Router
	logger         logger.Logger
	authMiddleware *auth.MiddlewareAuth
}

func New(logger logger.Logger, tokenRepository cache.TokenRepositoryPort, contextPath string) (*HttpRouter, error) {
	normalizedPath, err := normalizeContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	router := routerAdapter.NewRouter()

	authMiddleware := auth.NewAuthMiddleware(tokenRepository, logger)

	return &HttpRouter{
		ContextPath:    normalizedPath,
		Router:         router,
		logger:         logger,
		authMiddleware: authMiddleware,
	}, nil
}

func (r *HttpRouter) GetRouter() portRouter.Router {
	return r.Router
}

func normalizeContextPath(raw string) (string, error) {
	if raw == "" {
		return "", fmt.Errorf(constants.ErrorContextPathEmpty)
	}

	if strings.Contains(raw[1:], "/") {
		return "", fmt.Errorf(constants.ErrorContextPathSlash)
	}

	if raw[0] != '/' {
		raw = "/" + raw
	}

	if len(raw) <= 1 {
		return "", fmt.Errorf(constants.InvalidContextPath, raw)
	}

	return raw, nil
}
