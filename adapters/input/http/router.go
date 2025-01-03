package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/middlewares"
	"go.uber.org/zap"
)

type Router struct {
	ContextPath string
	chiRouter   chi.Router
	LoggerSugar *zap.SugaredLogger
}

func GetNewRouter(loggerSugar *zap.SugaredLogger) Router {
	router := chi.NewRouter()
	return Router{
		chiRouter:   router,
		LoggerSugar: loggerSugar,
	}
}

func (router Router) GetChiRouter() chi.Router {
	return router.chiRouter
}

func (router Router) AddGroupHandlerHealthCheck(ah *handlers.Generic) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/health-check", func(r chi.Router) {
			r.Get("/", ah.HealthCheck)
		})
	}
}

func (router Router) AddGroupHandlerUser(ah *handlers.User) func(r chi.Router) {
	return func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", ah.CreateUser)

			r.Group(func(r chi.Router) {
				r.Use(middlewares.AuthMiddleware(router.LoggerSugar))
				r.Get("/all", ah.GetAllUsers)
				r.Get("/{id}", ah.GetUserByID)
				r.Put("/{id}", ah.UpdateUser)
				r.Delete("/{id}", ah.SoftDeleteUser)
			})
		})
		r.Route("/login", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/", ah.Login)
			})
		})
	}
}
