package httpserver

import (
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
)

func ComposeRouter(deps *bootstrap.AppDependencies, contextPath string) (portRouter.Router, error) {
	httpRouter, err := NewHttpRouter(deps.Logger, deps.TokenRepository, contextPath)
	if err != nil {
		return nil, err
	}

	return BuildRouterRoutes(
		deps.Logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenRepository,
		httpRouter.BasePath,
		httpRouter.GetRouter(),
	)
}
