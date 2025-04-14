package server

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

func (r *RouteComposer) AddHealthCheckRoutes(gh *handlers.Generic) func(portRouter.Router) {
	return func(rt portRouter.Router) {
		rt.Route("/health-check", func(sub portRouter.Router) {
			sub.Get("/", gh.HealthCheckHandler)
		})
	}
}

func (r *RouteComposer) AddUserRoutes(uh *handlers.User) func(portRouter.Router) {
	return func(rt portRouter.Router) {
		rt.Route("/user", func(sub portRouter.Router) {
			sub.Post("/create", uh.CreateUserHandler)

			sub.Route("/", func(priv portRouter.Router) {
				priv.Use(r.authMiddleware.Auth)
				priv.Get("/all", uh.GetAllUsersHandler)
				priv.Get("/{id}", uh.GetUserByIDHandler)
				priv.Put("/", uh.UpdateUserHandler)
				priv.Put("/password", uh.UpdatePasswordHandler)
				priv.Delete("/", uh.SoftDeleteUserHandler)
			})
		})
	}
}

func (r *RouteComposer) AddAuthRoutes(ah *handlers.Auth) func(portRouter.Router) {
	return func(rt portRouter.Router) {
		rt.Route("/auth", func(priv portRouter.Router) {
			priv.Post("/login", ah.LoginHandler)

			priv.Route("/", func(protected portRouter.Router) {
				protected.Use(r.authMiddleware.Auth)
				protected.Post("/logout", ah.LogoutHandler)
			})
		})
	}
}
