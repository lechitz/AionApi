// Package http provides a set of utilities for building http servers.
package http

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// Params holds the parameters for building a new http.Server instance.
type Params struct {
	Handler           http.Handler
	Addr              string
	MaxHeaderBytes    int
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
}

// FromHTTP creates a new Params instance from the given http.Handler.
func FromHTTP(cfg *config.Config, h http.Handler) Params {
	return Params{
		Addr:              net.JoinHostPort(cfg.ServerHTTP.Host, cfg.ServerHTTP.Port),
		Handler:           h,
		ReadTimeout:       cfg.ServerHTTP.ReadTimeout,
		WriteTimeout:      cfg.ServerHTTP.WriteTimeout,
		ReadHeaderTimeout: cfg.ServerHTTP.ReadHeaderTimeout,
		IdleTimeout:       cfg.ServerHTTP.IdleTimeout,
		MaxHeaderBytes:    cfg.ServerHTTP.MaxHeaderBytes,
	}
}

// Build builds a new http.Server instance.
func Build(appCtx context.Context, p Params, _ logger.ContextLogger) *http.Server {
	return &http.Server{
		Addr:              p.Addr,
		Handler:           p.Handler,
		ReadTimeout:       p.ReadTimeout,
		WriteTimeout:      p.WriteTimeout,
		ReadHeaderTimeout: p.ReadHeaderTimeout,
		IdleTimeout:       p.IdleTimeout,
		MaxHeaderBytes:    p.MaxHeaderBytes,
		BaseContext:       func(_ net.Listener) context.Context { return appCtx },
	}
}
