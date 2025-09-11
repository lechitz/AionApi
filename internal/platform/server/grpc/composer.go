package grpc

import (
	"google.golang.org/grpc"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
)

// ComposeServer cria o *grpc.Server e (no futuro) registra os services.
func ComposeServer(_ interface{}, _ *bootstrap.AppDependencies, _ logger.ContextLogger) (*grpc.Server, error) {
	// TODO: implementar interceptors, registrations etc.
	return grpc.NewServer(), nil
}
