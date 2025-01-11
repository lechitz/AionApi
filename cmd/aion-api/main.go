package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	adpterHttpInput "github.com/lechitz/AionApi/adapters/input/http"
	"github.com/lechitz/AionApi/adapters/input/http/handlers"
	"github.com/lechitz/AionApi/adapters/output/cache"
	"github.com/lechitz/AionApi/adapters/output/db"
	"github.com/lechitz/AionApi/config"
	infraCache "github.com/lechitz/AionApi/infra/cache"
	infraDB "github.com/lechitz/AionApi/infra/db"
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

//var loggerSugar *zap.SugaredLogger

func InitLoggerSugar() (*zap.SugaredLogger, func()) {
	configZap := zap.NewProductionEncoderConfig()
	configZap.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(configZap)

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
	)

	loggerSugar := zap.New(core, zap.AddCaller())
	sugar := loggerSugar.Sugar()

	closeFunc := func() {
		if err := sugar.Sync(); err != nil {
			fmt.Printf("Failed to sync logger: %v\n", err)
		}
	}

	return sugar, closeFunc
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

	loggerSugar, closeLogger := InitLoggerSugar()
	defer closeLogger()

	// The GenerateJWTKey() function generates a unique secret key for JWT authentication and saves it to the .env file.
	////This function is commented out by default in the main() function.
	//utils.GenerateJWTKey()

	postgresConnectionDB := infraDB.NewPostgresDB(
		config.Setting.Postgres,
		loggerSugar,
	)
	defer infraDB.Close(postgresConnectionDB, loggerSugar)

	userPostgresDB := db.NewUserPostgresDB(
		postgresConnectionDB,
		loggerSugar,
	)

	redisClient := infraCache.NewRedisClient(
		config.Setting.RedisConfig.Addr,
		config.Setting.RedisConfig.Password,
		config.Setting.RedisConfig.DB,
		loggerSugar,
	)
	defer redisClient.Close()

	tokenStore := cache.NewTokenStore(
		redisClient,
		loggerSugar,
		config.Setting.SecretKey,
	)

	userService := &service.UserService{
		UserDomainDataBaseRepository: &userPostgresDB,
		UserDomainCacheRepository:    tokenStore,
		LoggerSugar:                  loggerSugar,
	}

	authService := &service.AuthService{
		UserDomainDataBaseRepository: &userPostgresDB,
		TokenStore:                   tokenStore,
		LoggerSugar:                  loggerSugar,
	}

	userHandler := &handlers.User{
		UserService: userService,
		LoggerSugar: loggerSugar,
	}

	authHandler := &handlers.Auth{
		AuthService: authService,
		LoggerSugar: loggerSugar,
	}

	genericHandler := &handlers.Generic{
		LoggerSugar: loggerSugar,
	}

	contextPath := config.Setting.Server.Context
	newRouter := adpterHttpInput.GetNewRouter(loggerSugar, tokenStore)
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
