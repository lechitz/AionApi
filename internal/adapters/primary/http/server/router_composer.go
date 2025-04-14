package server

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
)

func ComposeRouter(deps *bootstrap.AppDependencies, logger logger.Logger, contextPath string) (portRouter.Router, error) {
	httpRouter, err := NewHttpRouter(logger, deps.TokenRepository, contextPath)
	if err != nil {
		return nil, err
	}

	return BuildRouterRoutes(
		logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenRepository,
		httpRouter.BashPath,
		httpRouter.GetRouter(),
	)
}
