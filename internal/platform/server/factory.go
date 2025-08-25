package server

import (
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/config"
)

// FromHTTP creates server parameters for HTTP server from the provided configuration.
// It takes a configuration object and an HTTP handler as input and returns a Params struct
// containing all necessary server configuration parameters.
func FromHTTP(cfg *config.Config, handler http.Handler) Params {
	return Params{
		Name:              cfg.ServerHTTP.Name,
		Host:              cfg.ServerHTTP.Host,
		Port:              cfg.ServerHTTP.Port,
		Handler:           handler,
		ReadTimeout:       cfg.ServerHTTP.ReadTimeout,
		WriteTimeout:      cfg.ServerHTTP.WriteTimeout,
		ReadHeaderTimeout: cfg.ServerHTTP.ReadHeaderTimeout,
		IdleTimeout:       cfg.ServerHTTP.IdleTimeout,
		MaxHeaderBytes:    cfg.ServerHTTP.MaxHeaderBytes,
	}
}

// FromGraphQL initializes and returns Params for a GraphQL server using the given configuration and HTTP handler.
func FromGraphQL(cfg *config.Config, handler http.Handler) Params {
	return Params{
		Name:              cfg.ServerGraphql.Name,
		Host:              cfg.ServerGraphql.Host,
		Port:              cfg.ServerGraphql.Port,
		Handler:           handler,
		ReadTimeout:       cfg.ServerGraphql.ReadTimeout,
		WriteTimeout:      cfg.ServerGraphql.WriteTimeout,
		ReadHeaderTimeout: cfg.ServerGraphql.ReadHeaderTimeout,
		IdleTimeout:       cfg.ServerGraphql.IdleTimeout,
		MaxHeaderBytes:    cfg.ServerGraphql.MaxHeaderBytes,
	}
}
