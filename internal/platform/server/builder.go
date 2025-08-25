package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// serverStartedMessage is a constant string used to format the message indicating the start of the server.
const serverStartedMessage = "%s server started"

func Build(ctx context.Context, p Params, logger output.ContextLogger) *http.Server {
	addr := net.JoinHostPort(p.Host, p.Port)

	srv := &http.Server{
		Addr:              addr,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
		Handler:           p.Handler,
		ReadTimeout:       p.ReadTimeout,
		WriteTimeout:      p.WriteTimeout,
		ReadHeaderTimeout: p.ReadHeaderTimeout,
		IdleTimeout:       p.IdleTimeout,
		MaxHeaderBytes:    p.MaxHeaderBytes,
	}

	logger.Infow(fmt.Sprintf(serverStartedMessage, p.Name),
		commonkeys.ServerName, p.Name,
		commonkeys.ServerAddr, addr,
	)

	return srv
}
