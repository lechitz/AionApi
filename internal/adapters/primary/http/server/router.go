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

type RouteComposer struct {
	BashPath       string
	Router         portRouter.Router
	logger         logger.Logger
	authMiddleware *auth.MiddlewareAuth
}

func NewHttpRouter(logger logger.Logger, tokenRepository cache.TokenRepositoryPort, contextPath string) (*RouteComposer, error) {
	normalizedPath, err := normalizeContextPath(contextPath)
	if err != nil {
		return nil, err
	}

	return &RouteComposer{
		BashPath:       normalizedPath,
		Router:         routerAdapter.NewRouter(),
		logger:         logger,
		authMiddleware: auth.NewAuthMiddleware(tokenRepository, logger),
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
