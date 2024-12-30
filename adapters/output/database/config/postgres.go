package config

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewPostgresDB(DBUser, DBPassword, DBName, DBHost, DBPort string, loggerSugar *zap.SugaredLogger) *gorm.DB {
	var err error
	conString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost,
		DBPort,
		DBUser,
		DBPassword,
		DBName)

	loggerSugar.Infow("db connection", "connection_string", conString)

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
