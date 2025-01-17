package server

import (
	"fmt"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/app/bootstrap"
	"go.uber.org/zap"
)

//todo tá tudo certo.. erros estranhos, possíveis de serem do pacote router para o server, apos a mudanća de nome do pacote

func InitRouter(dependencies *bootstrap.AppDependencies, logger *zap.SugaredLogger, contextPath string) (*Router, error) {

	if contextPath == "" {
		return nil, fmt.Errorf("contextPath cannot be empty")
	}

	if strings.Contains(contextPath[1:], "/") {
		return nil, fmt.Errorf("contextPath cannot contain additional slashes `/`")
	}

	userHandler, authHandler, genericHandler := initializeHandlers(dependencies, logger)

	router, err := GetNewRouter(logger, dependencies.AuthService, dependencies.TokenService, contextPath)
	if err != nil {
		return nil, err
	}

	if err := configureRoutes(router, userHandler, authHandler, genericHandler); err != nil {
		return nil, err
	}

	return router, nil
}

func initializeHandlers(dependencies *bootstrap.AppDependencies, loggerSugar *zap.SugaredLogger) (*handlers.User, *handlers.Auth, *handlers.Generic) {
	return &handlers.User{
			UserService: dependencies.UserService,
			LoggerSugar: loggerSugar,
		},
		&handlers.Auth{
			AuthService: dependencies.AuthService,
			LoggerSugar: loggerSugar,
		},
		&handlers.Generic{
			LoggerSugar: loggerSugar,
		}
}

func configureRoutes(router *Router, userHandler *handlers.User, authHandler *handlers.Auth, genericHandler *handlers.Generic) error {

	contextPath := router.ContextPath

	if len(contextPath) < 1 || contextPath[0] != '/' {
		return fmt.Errorf("contextPath must start with a '/'")
	}

	router.GetChiRouter().Route(contextPath, func(r chi.Router) {
		r.NotFound(genericHandler.NotFoundHandler)
		r.Group(router.AddHealthCheckRoutes(genericHandler))
		r.Group(router.AddUserRoutes(userHandler))
		r.Group(router.AddAuthRoutes(authHandler))
	})

	return nil
}
