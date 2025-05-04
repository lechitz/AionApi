package httpserver

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/infra/bootstrap"
	"github.com/lechitz/AionApi/internal/infra/config"
	"net/http"
)

func NewHTTPServer(deps *bootstrap.AppDependencies, setting *config.Config) (*http.Server, error) {
	router, err := ComposeRouter(deps, setting.ServerHTTP.Context)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:           fmt.Sprintf(":%s", setting.ServerHTTP.Port),
		Handler:        router,
		ReadTimeout:    setting.ServerHTTP.ReadTimeout,
		WriteTimeout:   setting.ServerHTTP.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}, nil
}
