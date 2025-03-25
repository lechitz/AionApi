package server

import (
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	inputHttp "github.com/lechitz/AionApi/ports/input/http"
	outputHttp "github.com/lechitz/AionApi/ports/output/security"

	"go.uber.org/zap"
)

func InitRouter(
	logger *zap.SugaredLogger,
	userService inputHttp.IUserService,
	authService inputHttp.IAuthService,
	tokenService outputHttp.ITokenService,
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
