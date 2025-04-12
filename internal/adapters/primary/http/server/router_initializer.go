package server

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/server/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"strings"

	"github.com/go-chi/chi/v5"
	inputHttp "github.com/lechitz/AionApi/internal/core/ports/input/http"
)

func InitRouter(logger logger.Logger, userService inputHttp.UserService, authService inputHttp.AuthService, tokenRepository cache.TokenRepositoryPort, contextPath string) (*Router, error) {
	if contextPath == "" {
		return nil, fmt.Errorf(constants.ErrorContextPathEmpty)
	}
	if strings.Contains(contextPath[1:], "/") {
		return nil, fmt.Errorf(constants.ErrorContextPathSlash)
	}

	userHandler := &handlers.User{
		UserService: userService,
		Logger:      logger,
	}
	authHandler := &handlers.Auth{
		AuthService: authService,
		Logger:      logger,
	}
	genericHandler := &handlers.Generic{
		Logger: logger,
	}

	router, err := GetNewRouter(logger, tokenRepository, contextPath)
	if err != nil {
		return nil, err
	}

	err = configureRoutes(router, userHandler, authHandler, genericHandler)
	if err != nil {
		return nil, err
	}

	return router, nil
}

func configureRoutes(router *Router, uh *handlers.User, ah *handlers.Auth, gh *handlers.Generic) error {
	contextPath := router.ContextPath

	if len(contextPath) < 1 || contextPath[0] != '/' {
		return fmt.Errorf(constants.ErrorContextPathSlash)
	}

	router.GetChiRouter().Route(contextPath, func(r chi.Router) {
		r.NotFound(gh.NotFoundHandler)
		r.Group(router.AddHealthCheckRoutes(gh))
		r.Group(router.AddUserRoutes(uh))
		r.Group(router.AddAuthRoutes(ah))
	})

	return nil
}
