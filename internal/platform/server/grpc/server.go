package grpc

import (
	"github.com/lechitz/AionApi/internal/platform/config"
)

// Params is the parameters for the server.
type Params struct {
	Addr string
}

// FromGRPC creates a new Params instance from the given config.
func FromGRPC(cfg *config.Config) Params {
	return Params{Addr: cfg.ServerGRPC.Addr}
}
