// Package httpserver provides a wrapper around chi.Router implementing HTTPRouter.
package httpserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router wraps a chi.Router to implement HTTPRouter.
type Router struct {
	chi          chi.Router
	errorHandler func(http.ResponseWriter, *http.Request, error)
}

// NewRouter initializes and returns a new HTTPRouter.
func NewRouter() HTTPRouter {
	return &Router{chi: chi.NewRouter()}
}

func (c *Router) Use(middleware func(http.Handler) http.Handler) { c.chi.Use(middleware) }

func (c *Router) Route(pattern string, fn func(r HTTPRouter)) {
	sub := chi.NewRouter()
	wrapped := &Router{chi: sub}
	fn(wrapped)
	c.chi.Mount(pattern, sub)
}

func (c *Router) Get(path string, handler http.HandlerFunc)        { c.chi.Get(path, handler) }
func (c *Router) Post(path string, handler http.HandlerFunc)       { c.chi.Post(path, handler) }
func (c *Router) Put(path string, handler http.HandlerFunc)        { c.chi.Put(path, handler) }
func (c *Router) Delete(path string, handler http.HandlerFunc)     { c.chi.Delete(path, handler) }
func (c *Router) Mount(pattern string, handler http.Handler)       { c.chi.Mount(pattern, handler) }
func (c *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) { c.chi.ServeHTTP(w, r) }
func (c *Router) Group(fn func(r HTTPRouter))                      { fn(c) }

func (c *Router) SetNotFoundHandler(handler http.HandlerFunc) { c.chi.NotFound(handler) }
func (c *Router) SetMethodNotAllowedHandler(handler http.HandlerFunc) {
	c.chi.MethodNotAllowed(handler)
}
func (c *Router) SetErrorHandler(handler func(http.ResponseWriter, *http.Request, error)) {
	c.errorHandler = handler
}

func (c *Router) GroupWithMiddleware(middleware func(http.Handler) http.Handler, fn func(r HTTPRouter)) {
	c.chi.Group(func(r chi.Router) {
		r.Use(middleware)
		wrapped := &Router{chi: r}
		fn(wrapped)
	})
}
