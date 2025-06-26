// Package router defines an interface for managing HTTP routes.
package router

import "net/http"

// Router defines an interface to manage and configure HTTP routes.
// Use adds middleware functions to process HTTP requests.
// Route initializes sub-routes for specific patterns within the current router.
// Get registers a handler for the GET HTTP method at a specified path.
// Post registers a handler for the POST HTTP method at a specified path.
// Put registers a handler for the PUT HTTP method at a specified path.
// Delete registers a handler for the DELETE HTTP method at a specified path.
// Mount attaches an external HTTP handler to a specified pattern.
// ServeHTTP enables the router to fulfill the http.Handler interface.
// Group organizes routes and shared middlewares under a logical grouping.
type Router interface {
	Use(middleware func(http.Handler) http.Handler)

	Route(pattern string, fn func(r Router))

	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)

	Mount(pattern string, handler http.Handler)

	ServeHTTP(w http.ResponseWriter, r *http.Request)

	Group(fn func(r Router))
}
