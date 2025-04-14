package server

import (
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
)

func BuildRouter(deps *bootstrap.AppDependencies, logger logger.Logger, contextPath string) (portRouter.Router, error) {
	httpRouter, err := New(logger, deps.TokenRepository, contextPath)
	if err != nil {
		return nil, err
	}

	builder, err := InitRouter(
		logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenRepository,
		contextPath,
		httpRouter.GetRouter(),
	)
	if err != nil {
		return nil, err
	}

	return builder.Router, nil
}
