package server

import (
	"net/http"
	"time"
)

// Params gathers all necessary parameters to build an http.Server.
type Params struct {
	Name              string
	Host              string
	Port              string
	Handler           http.Handler
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	ReadHeaderTimeout time.Duration
	IdleTimeout       time.Duration
	MaxHeaderBytes    int
}
