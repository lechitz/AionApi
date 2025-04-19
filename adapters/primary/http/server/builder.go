package server

import (
	"github.com/lechitz/AionApi/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/adapters/primary/http/middleware/auth"
	contextbuilder "github.com/lechitz/AionApi/adapters/primary/http/middleware/contextbuilder"
	"github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

type RouterBuilder struct {
	ContextPath string
	Router      portRouter.Router
}

func BuildRouterRoutes(logger logger.Logger, userService http.UserService, authService http.AuthService, tokenRepo cache.TokenRepositoryPort, contextPath string, adapter portRouter.Router) (portRouter.Router, error) {

	adapter.Use(contextbuilder.InjectRequestIDMiddleware)

	genericHandler := handlers.NewGeneric(logger)
	userHandler := handlers.NewUser(userService, logger)
	authHandler := handlers.NewAuth(authService, logger)

	authMiddleware := auth.NewAuthMiddleware(tokenRepo, logger)

	r := &RouteComposer{
		BasePath:       contextPath,
		Router:         adapter,
		logger:         logger,
		authMiddleware: authMiddleware,
	}

	adapter.Route(contextPath, func(rt portRouter.Router) {
		rt.Group(r.AddHealthCheckRoutes(genericHandler))
		rt.Group(r.AddUserRoutes(userHandler))
		rt.Group(r.AddAuthRoutes(authHandler))
	})

	return adapter, nil
}
