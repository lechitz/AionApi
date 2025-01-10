package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/middlewares"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath    string
	chiRouter      chi.Router
	LoggerSugar    *zap.SugaredLogger
	AuthMiddleware *middlewares.AuthMiddleware
}

func GetNewRouter(loggerSugar *zap.SugaredLogger, tokenStore output.ITokenStore) Router {
	router := chi.NewRouter()

	authMiddleware := middlewares.NewAuthMiddleware(tokenStore, loggerSugar)

	return Router{
		chiRouter:      router,
		LoggerSugar:    loggerSugar,
		AuthMiddleware: authMiddleware,
	}
}

func (router Router) GetChiRouter() chi.Router {
	return router.chiRouter
}

func (router Router) AddGroupHandlerHealthCheck(ah *handlers.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", ah.HealthCheckHandler)
		})
	}
}

func (router Router) AddGroupHandlerUser(ah *handlers.User) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", ah.CreateUserHandler)

			r.Group(func(r chi.Router) {
				r.Use(router.AuthMiddleware.Auth)
				r.Get("/all", ah.GetAllUsersHandler)
				r.Get("/{id}", ah.GetUserByIDHandler)
				r.Put("/{id}", ah.UpdateUserHandler)
				r.Delete("/{id}", ah.SoftDeleteUserHandler)
			})
		})
	}
}

func (router Router) AddGroupHandlerAuth(ah *handlers.Auth) func(r chi.Router) {
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
