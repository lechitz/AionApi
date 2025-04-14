package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	portRouter "github.com/lechitz/AionApi/internal/core/ports/output/router"
)

type ChiRouter struct {
	chi chi.Router
}

func NewRouter() portRouter.Router {
	return &ChiRouter{chi: chi.NewRouter()}
}

func (c *ChiRouter) Use(middleware func(http.Handler) http.Handler) {
	c.chi.Use(middleware)
}

func (c *ChiRouter) Route(pattern string, fn func(r portRouter.Router)) {
	c.chi.Route(pattern, func(r chi.Router) {
		fn(&ChiRouter{chi: r})
	})
}

func (c *ChiRouter) Get(path string, handler http.HandlerFunc) {
	c.chi.Get(path, handler)
}

func (c *ChiRouter) Post(path string, handler http.HandlerFunc) {
	c.chi.Post(path, handler)
}

func (c *ChiRouter) Put(path string, handler http.HandlerFunc) {
	c.chi.Put(path, handler)
}

func (c *ChiRouter) Delete(path string, handler http.HandlerFunc) {
	c.chi.Delete(path, handler)
}

func (c *ChiRouter) Mount(pattern string, handler http.Handler) {
	c.chi.Mount(pattern, handler)
}

func (c *ChiRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.chi.ServeHTTP(w, r)
}

func (c *ChiRouter) Group(fn func(r portRouter.Router)) {
	c.chi.Group(func(r chi.Router) {
		fn(&ChiRouter{chi: r})
	})
}
