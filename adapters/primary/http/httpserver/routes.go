package httpserver

import (
	"github.com/lechitz/AionApi/adapters/primary/http/handlers"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

// AddHealthCheckRoutes registers the health check route and its handler to the provided router.
func (r *RouteComposer) AddHealthCheckRoutes(gh *handlers.Generic) func(portRouter.Router) {
	return func(healthGroup portRouter.Router) {
		healthGroup.Route("/health-check", func(healthProtected portRouter.Router) {
			healthProtected.Get("/", gh.HealthCheckHandler)
		})
	}
}

// AddUserRoutes registers user-related routes with respective handlers and applies authentication middleware to protected routes.
func (r *RouteComposer) AddUserRoutes(uh *handlers.User) func(portRouter.Router) {
	return func(rt portRouter.Router) {
		rt.Route("/user", func(userGroup portRouter.Router) {
			userGroup.Post("/create", uh.CreateUserHandler)

			userGroup.Route("/", func(userProtected portRouter.Router) {
				userProtected.Use(r.authMiddleware.Auth)

				userProtected.Get("/all", uh.GetAllUsersHandler)
				userProtected.Get("/{user_id}", uh.GetUserByIDHandler)
				userProtected.Put("/", uh.UpdateUserHandler)
				userProtected.Put("/password", uh.UpdatePasswordHandler)
				userProtected.Delete("/", uh.SoftDeleteUserHandler)
			})
		})
	}
}

// AddAuthRoutes registers authentication routes including login and logout endpoints, applying authentication middleware where required.
func (r *RouteComposer) AddAuthRoutes(ah *handlers.Auth) func(portRouter.Router) {
	return func(rt portRouter.Router) {
		rt.Route("/auth", func(authGroup portRouter.Router) {
			authGroup.Post("/login", ah.LoginHandler)

			authGroup.Route("/", func(authProtected portRouter.Router) {
				authProtected.Use(r.authMiddleware.Auth)

				authProtected.Post("/logout", ah.LogoutHandler)
			})
		})
	}
}
