package router

import "net/http"

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
