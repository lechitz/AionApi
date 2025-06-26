package httpserver

import (
	"fmt"
	"net/http"

	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// NewHTTPServer creates and configures a new HTTP server using provided dependencies and configuration. Returns the server instance or an error.
func NewHTTPServer(deps *bootstrap.AppDependencies, setting *config.Config) (*http.Server, error) {
	router, err := ComposeRouter(deps, setting.ServerHTTP.Context)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.ServerHTTP.Port),
		Handler:        otelhttp.NewHandler(router, "AionApi-REST"),
		ReadTimeout:    setting.ServerHTTP.ReadTimeout,
		WriteTimeout:   setting.ServerHTTP.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}, nil
}
