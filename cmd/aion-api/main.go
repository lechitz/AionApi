package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/lechitz/AionApi/adapters/input/http"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/output/db"
	"github.com/lechitz/AionApi/adapters/output/repository"
	"github.com/lechitz/AionApi/config"
	"github.com/lechitz/AionApi/internal/core/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

const (
	ServerStarted                = "server started"
	ErrorToListenAndStartsServer = "error to listen and starts server"
)

var loggerSugar *zap.SugaredLogger

func init() {
	err := envconfig.Process("setting", &config.Setting)
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
	//middlewares.GenerateJWTKey()
}

func main() {
	postgresConnectionDB := db.NewPostgresDB(
		config.Setting.Postgres.DBUser,
		config.Setting.Postgres.DBPassword,
		config.Setting.Postgres.DBName,
		config.Setting.Postgres.DBHost,
		config.Setting.Postgres.DBPort,
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

	loginPostgresDB := repository.NewLoginPostgresDB(postgresConnectionDB, loggerSugar)

	loginService := &service.LoginService{
		LoginDomainDataBaseRepository: &loginPostgresDB,
		LoggerSugar:                   loggerSugar,
	}

	loginHandler := &handlers.Login{
		LoginService: loginService,
		LoggerSugar:  loggerSugar,
	}

	genericHandler := &handlers.Generic{
		LoggerSugar: loggerSugar,
	}

	contextPath := config.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar)
	newRouter.GetChiRouter().Route(fmt.Sprintf("/%s", contextPath), func(r chi.Router) {
		r.NotFound(genericHandler.NotFound)
		r.Group(newRouter.AddGroupHandlerHealthCheck(genericHandler))
		r.Group(newRouter.AddGroupHandlerUser(userHandler))
		r.Group(newRouter.AddGroupHandlerLogin(loginHandler))
	})

	serverHttp := &http.Server{
		Addr:           fmt.Sprintf(":%s", config.Setting.Server.Port),
		Handler:        newRouter.GetChiRouter(),
		ReadTimeout:    config.Setting.Server.ReadTimeout,
		WriteTimeout:   config.Setting.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	loggerSugar.Infow(ServerStarted, "port", serverHttp.Addr, "contextPath", contextPath)

	if err := serverHttp.ListenAndServe(); err != nil {
		loggerSugar.Errorw(ErrorToListenAndStartsServer, "port", serverHttp.Addr,
			"contextPath", contextPath, "err", err.Error())
		panic(err.Error())
	}
}
