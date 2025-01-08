package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/lechitz/AionApi/adapters/input/http"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/output/repository"
	"github.com/lechitz/AionApi/config"
	"github.com/lechitz/AionApi/infra/db"
	"github.com/lechitz/AionApi/internal/core/service"
	"github.com/lechitz/AionApi/pkg/utils"
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

func InitLogger() {
	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(configZap)
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)
	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	loggerSugar = logger.Sugar()
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	err = envconfig.Process("setting", &config.Setting)
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	LoadConfig()
	InitLogger()
	utils.GenerateJWTKey()

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

	authService := &service.AuthService{
		UserDomainDataBaseRepository: &userPostgresDB,
		LoggerSugar:                  loggerSugar,
	}

	authHandler := &handlers.Auth{
		AuthService: authService,
		LoggerSugar: loggerSugar,
	}

	genericHandler := &handlers.Generic{
		LoggerSugar: loggerSugar,
	}

	contextPath := config.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar)
	newRouter.GetChiRouter().Route(fmt.Sprintf("/%s", contextPath), func(r chi.Router) {
		r.NotFound(genericHandler.NotFoundHandler)
		r.Group(newRouter.AddGroupHandlerHealthCheck(genericHandler))
		r.Group(newRouter.AddGroupHandlerUser(userHandler))
		r.Group(newRouter.AddGroupHandlerAuth(authHandler))
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
		loggerSugar.Errorw(ErrorToListenAndStartsServer, "port", serverHttp.Addr, "contextPath", contextPath, "err", err.Error())
		panic(err.Error())
	}
}
