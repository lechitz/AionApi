package middlewares

import (
	"fmt"
	"github.com/lechitz/AionApi/src/auth"
	"github.com/lechitz/AionApi/src/responses"
	"log"
	"net/http"
)

func Logger(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Logging...")
		log.Printf("\n%s - %s %s\n", r.Method, r.RequestURI, r.Host)
		nextFunc(w, r)
	}
}

func Authentication(nextFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := auth.ValidateToken(r); err != nil {
			responses.Error(w, http.StatusUnauthorized, err)
			return
		}

		nextFunc(w, r)
	}
}
