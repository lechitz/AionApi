package server

import (
	"fmt"
	"net/http"

	"github.com/lechitz/AionApi/app/bootstrap"
	"github.com/lechitz/AionApi/app/config"
	"go.uber.org/zap"
)

// NewHTTPServer creates and configures an HTTP server.
func NewHTTPServer(deps *bootstrap.AppDependencies, logger *zap.SugaredLogger, setting *config.Config) (*http.Server, error) {

	route, err := InitRouter(deps, logger, setting.Server.Context)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.Server.Port),
		Handler:        route.GetChiRouter(),
		ReadTimeout:    setting.Server.ReadTimeout,
		WriteTimeout:   setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	return srv, nil
}
