package server

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"net/http"
)

func NewHTTPServer(deps *bootstrap.AppDependencies, logger logger.Logger, setting *config.Config) (*http.Server, error) {
	router, err := ComposeRouter(deps, logger, setting.Server.Context)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.Server.Port),
		Handler:        router,
		ReadTimeout:    setting.Server.ReadTimeout,
		WriteTimeout:   setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}, nil
}
