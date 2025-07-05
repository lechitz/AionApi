// Package chi provides a wrapper around the chi.Router to implement the output.Router interface.
package chi

import (
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router is a wrapper around the chi.Router to implement the output.Router interface.
type Router struct {
	chi chi.Router
}

// NewRouter initializes and returns a new instance of output.Router using a Router implementation.
func NewRouter() output.Router {
	return &Router{chi: chi.NewRouter()}
}

// Use adds a middleware function to the router's middleware stack.
func (c *Router) Use(middleware func(http.Handler) http.Handler) {
	c.chi.Use(middleware)
}

// Route defines a sub-router for a specific route pattern within the current router.
func (c *Router) Route(pattern string, fn func(r output.Router)) {
	c.chi.Route(pattern, func(r chi.Router) {
		fn(&Router{chi: r})
	})
}

// Get registers a route that matches GET HTTP method for the specified path.
func (c *Router) Get(path string, handler http.HandlerFunc) {
	c.chi.Get(path, handler)
}

// Post registers a route that matches the POST HTTP method for the specified path with the given handler.
func (c *Router) Post(path string, handler http.HandlerFunc) {
	c.chi.Post(path, handler)
}

// Put registers a route that matches the PUT HTTP method for the specified path with the provided handler function.
func (c *Router) Put(path string, handler http.HandlerFunc) {
	c.chi.Put(path, handler)
}

// Delete registers a route that matches the DELETE HTTP method for the specified path with the provided handler function.
func (c *Router) Delete(path string, handler http.HandlerFunc) {
	c.chi.Delete(path, handler)
}

// Mount attaches another HTTP handler to the specified URL pattern within the router.
func (c *Router) Mount(pattern string, handler http.Handler) {
	c.chi.Mount(pattern, handler)
}

// ServeHTTP delegates the HTTP request and response handling to the underlying chi.Router.
func (c *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.chi.ServeHTTP(w, r)
}

// Group creates a new router group where shared middlewares or routes can be defined and applied.
func (c *Router) Group(fn func(r output.Router)) {
	c.chi.Group(func(r chi.Router) {
		fn(&Router{chi: r})
	})
}
