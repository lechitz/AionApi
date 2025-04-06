package db

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/infrastructure/db/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewDatabaseConnection(cfg config.DBConfig, loggerSugar *zap.SugaredLogger) *gorm.DB {
	var err error
	conString := fmt.Sprintf(constants.MsgFormatConString, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	loggerSugar.Infow(constants.MsgDBConnection, constants.Host, cfg.DBHost, constants.Port, cfg.DBPort, constants.DBName, cfg.DBName)

	DB, err := connecting(conString, loggerSugar)
	if err != nil {
		log.Panic(err)
	}

	return DB
}

func connecting(conString string, loggerSugar *zap.SugaredLogger) (*gorm.DB, error) {

	tryConnect := 1

	for {
		loggerSugar.Infow(constants.MsgTryingStartsPostgresDB, "try", tryConnect)
		DB, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err != nil && tryConnect != 3 {

			tryConnect++
			if tryConnect > 3 {
				loggerSugar.Infow(constants.ErrorToStartsThePostgresDB, constants.MsgTryingToConnect, tryConnect)
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
		loggerSugar.Errorw(constants.MsgToRetrieveSQLFromGorm, constants.Error, err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		loggerSugar.Errorw(constants.ErrorToCloseThePostgresDB, constants.Error, err.Error())
	}
}
