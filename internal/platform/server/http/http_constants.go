// Package http centralizes HTTP server constants to avoid magic strings.
// NOTE: Values here serve as defaults. In the composer, we allow overrides via cfg.
//
//revive:disable:var-naming // keep package name aligned with HTTP server module naming
package http

//revive:enable:var-naming

// Default routes and mount points.
const (
	// DefaultRouteHealth Default health-check route inside the API context (e.g., /aion/health).
	DefaultRouteHealth = "/health"

	// DefaultSwaggerMountPath Where Swagger UI+assets are mounted under the API context (e.g., /aion/swagger/*).
	DefaultSwaggerMountPath = "/swagger"

	// DefaultDocsAliasPath Convenience alias that redirects to {SwaggerMountPath}/index.html (e.g., /aion/docs).
	DefaultDocsAliasPath = "/docs"
)

// Swagger/OpenAPI related defaults.
const (
	// DefaultSwaggerDocFile File names used by the UI; combined with DefaultSwaggerMountPath.
	DefaultSwaggerDocFile   = "doc.json"
	DefaultSwaggerIndexFile = "index.html"
)

// Logging messages.
const (
	LogErrComposeGraphQL      = "failed to compose GraphQL handler"
	OTelHTTPHandlerNameFormat = "%s-HTTP"
)
