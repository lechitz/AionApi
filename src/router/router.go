package router

import (
	"github.com/gorilla/mux"
	"github.com/lechitz/AionApi/src/router/rotas"
)

// GenareteRouter is a function that returns a new router
func GenareteRouter() *mux.Router {
	router := mux.NewRouter()
	return rotas.Configure(router)
}
