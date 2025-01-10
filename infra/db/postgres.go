package db

import (
	"fmt"
	"github.com/lechitz/AionApi/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewPostgresDB(config config.PostgresConfig, loggerSugar *zap.SugaredLogger) *gorm.DB {
	var err error
	conString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	loggerSugar.Infow("db connection", "host", config.DBHost, "port", config.DBPort, "dbname", config.DBName)

	DB, err := connecting(conString, loggerSugar)
	if err != nil {
		log.Panic(err)
	}

	return DB
}

func connecting(conString string, loggerSugar *zap.SugaredLogger) (*gorm.DB, error) {

	tryConnect := 1

	for {
		loggerSugar.Infow("trying starts postgres db", "try", tryConnect)
		DB, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err != nil && tryConnect != 3 {

			tryConnect++
			if tryConnect > 3 {
				loggerSugar.Infow("error to starts the postgres db", "tries to starts", tryConnect)
				return nil, err
			}

			time.Sleep(3 * time.Second)
			continue
		}

		return DB, err
	}
}

func Close(DB *gorm.DB, loggerSugar *zap.SugaredLogger) {
	sqlDB, err := DB.DB()
	if err != nil {
		loggerSugar.Errorw("Failed to retrieve SQL DB from Gorm", "error", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		loggerSugar.Errorw("Failed to close Postgres connection", "error", err.Error())
	}
}
