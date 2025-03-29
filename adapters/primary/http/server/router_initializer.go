package server

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/core/ports/input/http"
	"github.com/lechitz/AionApi/internal/core/ports/output/security"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/primary/http/handlers"
	"go.uber.org/zap"
)

func InitRouter(
	logger *zap.SugaredLogger,
	userService http.IUserService,
	authService http.IAuthService,
	tokenService security.ITokenService,
	contextPath string,
) (*Router, error) {

	if contextPath == "" {
		return nil, fmt.Errorf(ErrorContextPathEmpty)
	}
	if strings.Contains(contextPath[1:], "/") {
		return nil, fmt.Errorf(ErrorContextPathSlash)
	}

	userHandler := &handlers.User{
		UserService: userService,
		LoggerSugar: logger,
	}
	authHandler := &handlers.Auth{
		AuthService: authService,
		LoggerSugar: logger,
	}
	genericHandler := &handlers.Generic{
		LoggerSugar: logger,
	}

	router, err := GetNewRouter(logger, authService, tokenService, contextPath)
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
		return fmt.Errorf(ErrorContextPathSlash)
	}

	router.GetChiRouter().Route(contextPath, func(r chi.Router) {
		r.NotFound(gh.NotFoundHandler)
		r.Group(router.AddHealthCheckRoutes(gh))
		r.Group(router.AddUserRoutes(uh))
		r.Group(router.AddAuthRoutes(ah))
	})

	return nil
}
