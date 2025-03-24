package server

import (
	"fmt"
	"github.com/lechitz/AionApi/core/service"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/app/middleware/auth"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath    string
	Router         chi.Router
	LoggerSugar    *zap.SugaredLogger
	AuthMiddleware *auth.MiddlewareAuth
}

func GetNewRouter(loggerSugar *zap.SugaredLogger, authService *service.AuthService, tokenService *service.TokenService, contextPath string) (*Router, error) {

	if len(contextPath) > 0 && contextPath[0] != '/' {
		contextPath = "/" + contextPath
	}

	if len(contextPath) <= 1 {
		return nil, fmt.Errorf("invalid context path: '%s'", contextPath)
	}

	r := chi.NewRouter()

	authMiddleware := auth.NewAuthMiddleware(authService, tokenService, loggerSugar)

	router := &Router{
		ContextPath:    contextPath,
		Router:         r,
		LoggerSugar:    loggerSugar,
		AuthMiddleware: authMiddleware,
	}

	return router, nil
}

func (router *Router) GetChiRouter() chi.Router {
	return router.Router
}

func (router *Router) AddHealthCheckRoutes(ah *handlers.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", ah.HealthCheckHandler)
		})
	}
}

func (router *Router) AddUserRoutes(uh *handlers.User) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", uh.CreateUserHandler)

			r.Group(func(r chi.Router) {
				r.Use(router.AuthMiddleware.Auth)

				r.Get("/all", uh.GetAllUsersHandler)
				r.Get("/{id}", uh.GetUserByIDHandler)
				r.Put("/{id}", uh.UpdateUserHandler)
				r.Put("/password/{id}", uh.UpdatePasswordHandler)
				r.Delete("/{id}", uh.SoftDeleteUserHandler)
			})
		})
	}
}

func (router *Router) AddAuthRoutes(ah *handlers.Auth) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", ah.LoginHandler)

			r.Group(func(r chi.Router) {
				r.Use(router.AuthMiddleware.Auth)
				r.Post("/logout/{id}", ah.LogoutHandler)
			})
		})
	}
}
