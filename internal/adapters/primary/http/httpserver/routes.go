package httpserver

import (
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/auth"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/generic"
	"github.com/lechitz/AionApi/internal/adapters/primary/http/handlers/user"
	"github.com/lechitz/AionApi/internal/core/ports/output"
)

// AddHealthCheckRoutes registers the health check route and its handler to the provided router.
func (r *RouteComposer) AddHealthCheckRoutes(gh *generic.Handler) func(output.Router) {
	return func(healthGroup output.Router) {
		healthGroup.Route("/health-check", func(healthProtected output.Router) {
			healthProtected.Get("/", gh.HealthCheck)
		})
	}
}

// AddUserRoutes registers user-related routes with respective handlers and applies authentication middleware to protected routes.
func (r *RouteComposer) AddUserRoutes(uh *user.Handler) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/user", func(userGroup output.Router) {
			userGroup.Post("/create", uh.Create)

			userGroup.Route("/", func(userProtected output.Router) {
				userProtected.Use(r.authMiddleware.Auth)

				userProtected.Get("/all", uh.GetAllUsers)
				userProtected.Get("/{user_id}", uh.GetUserByID)
				userProtected.Put("/", uh.UpdateUser)
				userProtected.Put("/password", uh.UpdateUserPassword)
				userProtected.Delete("/", uh.SoftDeleteUser)
			})
		})
	}
}

// AddAuthRoutes registers authentication routes including login and logout endpoints, applying authentication middleware where required.
func (r *RouteComposer) AddAuthRoutes(ah *auth.Handler) func(output.Router) {
	return func(rt output.Router) {
		rt.Route("/auth", func(authGroup output.Router) {
			authGroup.Post("/login", ah.Login)

			authGroup.Route("/", func(authProtected output.Router) {
				authProtected.Use(r.authMiddleware.Auth)

				authProtected.Post("/logout", ah.Logout)
			})
		})
	}
}
