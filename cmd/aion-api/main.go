package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/lechitz/AionApi/adapters/input/http"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/output/database/config"
	"github.com/lechitz/AionApi/adapters/output/database/repository"
	"github.com/lechitz/AionApi/config/environments"
	"github.com/lechitz/AionApi/internal/core/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var loggerSugar *zap.SugaredLogger

func init() {
	err := envconfig.Process("setting", &environments.Setting)
	if err != nil {
		panic(err.Error())
	}

	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(configZap)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any
	loggerSugar = logger.Sugar()

	//GenerateJWTKey is used to generate a new JWT key for the .env file one time
	//utils.GenerateJWTKey()
}

func main() {
	postgresConnectionDB := config.NewPostgresDB(
		environments.Setting.Postgres.DBUser,
		environments.Setting.Postgres.DBPassword,
		environments.Setting.Postgres.DBName,
		environments.Setting.Postgres.DBHost,
		environments.Setting.Postgres.DBPort,
		loggerSugar,
	)

	userPostgresDB := repository.NewUserPostgresDB(postgresConnectionDB, loggerSugar)

	userService := &service.UserService{
		UserDomainDataBaseRepository: &userPostgresDB,
		LoggerSugar:                  loggerSugar,
	}

	userHandler := &handlers.User{
		UserService: userService,
		LoggerSugar: loggerSugar,
	}

	genericHandler := &handlers.Generic{
		LoggerSugar: loggerSugar,
	}

	contextPath := environments.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar)
	newRouter.GetChiRouter().Route(fmt.Sprintf("/%s", contextPath), func(r chi.Router) {
		r.NotFound(genericHandler.NotFound)
		r.Group(newRouter.AddGroupHandlerHealthCheck(genericHandler))
		r.Group(newRouter.AddGroupHandlerUser(userHandler))
	})

	serverHttp := &http.Server{
		Addr:           fmt.Sprintf(":%s", environments.Setting.Server.Port),
		Handler:        newRouter.GetChiRouter(),
		ReadTimeout:    environments.Setting.Server.ReadTimeout,
		WriteTimeout:   environments.Setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	loggerSugar.Infow("server started", "port", serverHttp.Addr, "contextPath", contextPath)

	if err := serverHttp.ListenAndServe(); err != nil {
		loggerSugar.Errorw("error to listen and starts server", "port", serverHttp.Addr,
			"contextPath", contextPath, "err", err.Error())
		panic(err.Error())
	}
}
