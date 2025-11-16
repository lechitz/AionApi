// Package chi implements the Router port on top of chi v5.
package chi

import (
	"net/http"

	chiv5 "github.com/go-chi/chi/v5"
	"github.com/lechitz/AionApi/internal/platform/server/http/ports"
)

// router is the chi v5-backed implementation of the Router port.
type router struct {
	chi         chiv5.Router
	errorHandle func(http.ResponseWriter, *http.Request, error)
}

// New returns a chi v5-based Router with no middlewares by default.
// Platform-level middlewares (request-id, recovery, CORS, etc.) should be
// applied by the platform composer before registering domain modules.
func New() ports.Router {
	return &router{
		chi: chiv5.NewRouter(),
	}
}

// Use applies one or more middlewares to the chi pipeline.
// Nil middlewares are ignored defensively.
func (r *router) Use(mw ...ports.Middleware) {
	for _, m := range mw {
		if m != nil {
			r.chi.Use(m)
		}
	}
}

// Group creates a sub-route tree under 'prefix'.
// The provided function receives an isolated Router for that subtree.
func (r *router) Group(prefix string, fn func(ports.Router)) {
	r.chi.Route(prefix, func(cr chiv5.Router) {
		fn(&router{chi: cr})
	})
}

// GroupWith creates a sub-route tree and applies a middleware only to it.
func (r *router) GroupWith(m ports.Middleware, fn func(ports.Router)) {
	r.chi.Group(func(gr chiv5.Router) {
		if m != nil {
			gr.Use(m)
		}
		fn(&router{chi: gr})
	})
}

// Mount attaches a ready http.Handler under a prefix (e.g., /graphql).
func (r *router) Mount(prefix string, h http.Handler) {
	r.chi.Mount(prefix, h)
}

// Handle registers a route for an arbitrary HTTP method.
func (r *router) Handle(method, path string, h http.Handler) {
	r.chi.Method(method, path, h)
}

// Convenience helpers for common HTTP verbs.

// GET is a convenience method to register a route for the HTTP GET verb.
func (r *router) GET(path string, h http.Handler) { r.Handle(http.MethodGet, path, h) }

// POST is a convenience method to register a route for the HTTP POST verb.
func (r *router) POST(path string, h http.Handler) { r.Handle(http.MethodPost, path, h) }

// PUT is a convenience method to register a route for the HTTP PUT verb.
func (r *router) PUT(path string, h http.Handler) { r.Handle(http.MethodPut, path, h) }

// DELETE is a convenience method to register a route for the HTTP DELETE verb.
func (r *router) DELETE(path string, h http.Handler) { r.Handle(http.MethodDelete, path, h) }

// SetNotFound sets the custom 404 handler.
func (r *router) SetNotFound(h http.Handler) {
	r.chi.NotFound(h.ServeHTTP)
}

// SetMethodNotAllowed sets the custom 405 handler.
func (r *router) SetMethodNotAllowed(h http.Handler) {
	r.chi.MethodNotAllowed(h.ServeHTTP)
}

// SetError stores a centralized error handler (optional).
// Note: chi does not invoke this automatically; the platform may call it
// from an error/recovery middleware to standardize error responses.
func (r *router) SetError(h func(http.ResponseWriter, *http.Request, error)) {
	r.errorHandle = h
}

// ServeHTTP implements http.Handler and delegates to chi.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}
