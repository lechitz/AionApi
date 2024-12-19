package rotas

import (
	"github.com/gorilla/mux"
	"github.com/lechitz/AionApi/src/middlewares"
	"net/http"
)

type Route struct {
	URI          string
	Method       string
	Function     func(http.ResponseWriter, *http.Request)
	AuthRequired bool
}

func Configure(r *mux.Router) *mux.Router {
	routes := routerUsers
	routes = append(routes, routeLogin)

	for _, route := range routes {

		if route.AuthRequired {
			r.HandleFunc(route.URI, middlewares.Logger(middlewares.Authentication(route.Function))).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, middlewares.Logger(route.Function)).Methods(route.Method)
		}

	}
	return r
}
