// Package input router defines an interface for managing HTTP routes.
package input

import "net/http"

// HTTPRouter defines an interface to manage and configure HTTP routes.
type HTTPRouter interface {
	Use(middleware func(http.Handler) http.Handler)
	Route(pattern string, fn func(r HTTPRouter))
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
	Mount(pattern string, handler http.Handler)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Group(fn func(r HTTPRouter))
	SetNotFoundHandler(handler http.HandlerFunc)
	SetMethodNotAllowedHandler(handler http.HandlerFunc)
	SetErrorHandler(handler func(http.ResponseWriter, *http.Request, error))
	GroupWithMiddleware(middleware func(http.Handler) http.Handler, fn func(r HTTPRouter))
}
