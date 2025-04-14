package server

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

type RouterBuilder struct {
	ContextPath string
	Router      portRouter.Router
}

func InitRouter(logger logger.Logger, userService http.UserService, authService http.AuthService, tokenRepo cache.TokenRepositoryPort, contextPath string, adapter portRouter.Router) (*RouterBuilder, error) {

	genericHandler := handlers.NewGeneric(logger)
	userHandler := handlers.NewUser(userService, logger)
	authHandler := handlers.NewAuth(authService, logger)

	r := &HttpRouter{
		ContextPath:    contextPath,
		Router:         adapter,
		logger:         logger,
		authMiddleware: auth.NewAuthMiddleware(tokenRepo, logger),
	}

	adapter.Route(contextPath, func(rt portRouter.Router) {
		rt.Group(r.AddHealthCheckRoutes(genericHandler))
		rt.Group(r.AddUserRoutes(userHandler))
		rt.Group(r.AddAuthRoutes(authHandler))
	})

	return &RouterBuilder{
		ContextPath: contextPath,
		Router:      adapter,
	}, nil
}
