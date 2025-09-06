package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
)

// serverStartedMessage is a constant string used to format the message indicating the start of the server.
const serverStartedMessage = "%s server started"

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

func Build(ctx context.Context, p Params, logger logger.ContextLogger) *http.Server {
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
