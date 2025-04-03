package server

import (
	"fmt"
	"net/http"

	"github.com/lechitz/AionApi/config"
	"github.com/lechitz/AionApi/internal/platform/bootstrap"

	"go.uber.org/zap"
)

func NewHTTPServer(deps *bootstrap.AppDependencies, logger *zap.SugaredLogger, setting *config.Config) (*http.Server, error) {

	route, err := InitRouter(
		logger,
		deps.UserService,
		deps.AuthService,
		deps.TokenService,
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
