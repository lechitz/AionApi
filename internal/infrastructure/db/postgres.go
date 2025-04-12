package db

import (
	"fmt"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infrastructure/db/constants"
	"github.com/lechitz/AionApi/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func NewDatabaseConnection(cfg config.DBConfig, logger logger.Logger) *gorm.DB {
	var err error
	conString := fmt.Sprintf(constants.MsgFormatConString, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	logger.Infow(constants.MsgDBConnection, constants.Host, cfg.DBHost, constants.Port, cfg.DBPort, constants.DBName, cfg.DBName)

	DB, err := connecting(conString, logger)
	if err != nil {
		log.Panic(err)
	}

	return DB
}

func connecting(conString string, logger logger.Logger) (*gorm.DB, error) {

	tryConnect := 1

	for {
		logger.Infow(constants.MsgTryingStartsPostgresDB, "try", tryConnect)
		DB, err := gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err != nil && tryConnect != 3 {

			tryConnect++
			if tryConnect > 3 {
				logger.Infow(constants.ErrorToStartsThePostgresDB, constants.MsgTryingToConnect, tryConnect)
				return nil, err
			}

			time.Sleep(3 * time.Second)
			continue
		}

		return DB, err
	}
}

func Close(DB *gorm.DB, logger logger.Logger) {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Errorw(constants.MsgToRetrieveSQLFromGorm, constants.Error, err.Error())
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Errorw(constants.ErrorToCloseThePostgresDB, constants.Error, err.Error())
	}
}
