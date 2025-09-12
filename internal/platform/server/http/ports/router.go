// Package ports is the HTTP server ports interface.
package ports

import "net/http"

// Middleware is a function that wraps a http.Handler.
type Middleware func(http.Handler) http.Handler

// Router is the HTTP server router interface.
type Router interface {
	Use(mw ...Middleware)

	Group(prefix string, fn func(Router))
	GroupWith(mw Middleware, fn func(Router))
	Mount(prefix string, h http.Handler)

	Handle(method, path string, h http.Handler)
	GET(path string, h http.Handler)
	POST(path string, h http.Handler)
	PUT(path string, h http.Handler)
	DELETE(path string, h http.Handler)

	SetNotFound(h http.Handler)
	SetMethodNotAllowed(h http.Handler)
	SetError(func(http.ResponseWriter, *http.Request, error))

	ServeHTTP(http.ResponseWriter, *http.Request)
}
