package server

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"
	"github.com/lechitz/AionApi/internal/platform/config"
	"net/http"
)

func NewHTTPServer(deps *bootstrap.AppDependencies, logger logger.Logger, setting *config.Config) (*http.Server, error) {

	route, err := InitRouter(
		logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenRepository,
		setting.Server.Context,
	)
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
