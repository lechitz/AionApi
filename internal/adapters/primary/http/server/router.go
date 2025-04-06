package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/middleware/auth"
	"github.com/lechitz/AionApi/internal/core/ports/output/cache"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath    string
	Router         chi.Router
	LoggerSugar    *zap.SugaredLogger
	AuthMiddleware *auth.MiddlewareAuth
}

func GetNewRouter(loggerSugar *zap.SugaredLogger, tokenRepository cache.TokenRepositoryPort, contextPath string) (*Router, error) {
	if len(contextPath) > 0 && contextPath[0] != '/' {
		contextPath = "/" + contextPath
	}

	if len(contextPath) <= 1 {
		return nil, fmt.Errorf("invalid context path: '%s'", contextPath)
	}

	r := chi.NewRouter()
	authMiddleware := auth.NewAuthMiddleware(tokenRepository, loggerSugar)

	return &Router{
		ContextPath:    contextPath,
		Router:         r,
		LoggerSugar:    loggerSugar,
		AuthMiddleware: authMiddleware,
	}, nil
}

func (router *Router) GetChiRouter() chi.Router {
	return router.Router
}

func (router *Router) AddHealthCheckRoutes(gh *handlers.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", gh.HealthCheckHandler)
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
				r.Put("/", uh.UpdateUserHandler)
				r.Put("/password", uh.UpdatePasswordHandler)
				r.Delete("/", uh.SoftDeleteUserHandler)
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
				r.Post("/logout", ah.LogoutHandler)
			})
		})
	}
}
