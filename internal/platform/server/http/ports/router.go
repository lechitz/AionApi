package ports

import "net/http"

// Middleware de plataforma (chi/echo/gin convertem p/ isso).
type Middleware func(http.Handler) http.Handler

// Router é a ÚNICA superfície que os módulos conhecem.
type Router interface {
	// Middlewares
	Use(mw ...Middleware)

	// Composição
	Group(prefix string, fn func(Router))     // cria subárvore: r.Group("/user", fn)
	GroupWith(mw Middleware, fn func(Router)) // subárvore com mw aplicado
	Mount(prefix string, h http.Handler)      // pendura um http.Handler pronto (ex.: GraphQL)

	// Handlers HTTP
	Handle(method, path string, h http.Handler)
	GET(path string, h http.Handler)
	POST(path string, h http.Handler)
	PUT(path string, h http.Handler)
	DELETE(path string, h http.Handler)

	// Defaults / integração
	SetNotFound(h http.Handler)
	SetMethodNotAllowed(h http.Handler)
	SetError(func(http.ResponseWriter, *http.Request, error)) // opcional, p/ recovery

	ServeHTTP(http.ResponseWriter, *http.Request)
}
