package grpc

import (
	"github.com/lechitz/AionApi/internal/platform/config"
)

type Params struct {
	Addr string
}

func FromGRPC(cfg *config.Config) Params {
	// Ajuste para seus campos reais
	return Params{Addr: cfg.ServerGRPC.Addr}
}
