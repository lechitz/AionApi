package db

import (
	"fmt"
	"github.com/lechitz/AionApi/app/config"
	"github.com/lechitz/AionApi/pkg/contextkeys"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewDatabaseConnection(config config.DBConfig, loggerSugar *zap.SugaredLogger) *gorm.DB {
	var err error
	conString := fmt.Sprintf(msgFormatConString,
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	loggerSugar.Infow(msgDBConnection, contextkeys.Host, config.DBHost, contextkeys.Port, config.DBPort, contextkeys.DBName, config.DBName)

	DB, err := connecting(conString, loggerSugar)
	if err != nil {
		log.Panic(err)
	}

	return DB
}

func connecting(conString string, loggerSugar *zap.SugaredLogger) (*gorm.DB, error) {

	tryConnect := 1

	for {
		loggerSugar.Infow(msgTryingStartsPostgresDB, "try", tryConnect)
		DB, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err != nil && tryConnect != 3 {

			tryConnect++
			if tryConnect > 3 {
				loggerSugar.Infow(errorToStartsThePostgresDB, msgTryingToConnect, tryConnect)
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
		loggerSugar.Errorw(msgToRetrieveSQLFromGorm, "error", err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		loggerSugar.Errorw(errorToCloseThePostgresDB, "error", err.Error())
	}
}
