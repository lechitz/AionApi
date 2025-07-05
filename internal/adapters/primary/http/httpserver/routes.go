package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// AddHealthCheckRoutes registers the health check route and its handler to the provided router.
func (r *RouteComposer) AddHealthCheckRoutes(gh *handlers.Generic) func(output.Router) {
	return func(healthGroup output.Router) {
		healthGroup.Route("/health-check", func(healthProtected output.Router) {
			healthProtected.Get("/", gh.HealthCheckHandler)
		})
	}
}

// AddUserRoutes registers user-related routes with respective handlers and applies authentication middleware to protected routes.
func (r *RouteComposer) AddUserRoutes(uh *handlers.User) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/user", func(userGroup output.Router) {
			userGroup.Post("/create", uh.CreateUserHandler)

			userGroup.Route("/", func(userProtected output.Router) {
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
func (r *RouteComposer) AddAuthRoutes(ah *handlers.Auth) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/auth", func(authGroup output.Router) {
			authGroup.Post("/login", ah.LoginHandler)

			authGroup.Route("/", func(authProtected output.Router) {
				authProtected.Use(r.authMiddleware.Auth)

				authProtected.Post("/logout", ah.LogoutHandler)
			})
		})
	}
}
