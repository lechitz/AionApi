package db

import (
	"fmt"
	"github.com/lechitz/AionApi/adapters/secondary/db/postgres/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output/logger"
	"github.com/lechitz/AionApi/internal/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func NewDatabaseConnection(cfg config.DBConfig, logger logger.Logger) (*gorm.DB, error) {
	conString := fmt.Sprintf(constants.MsgFormatConString, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	logger.Infow(constants.MsgDBConnection, constants.Host, cfg.Host, constants.Port, cfg.Port, constants.DBName, cfg.Name)

	db, err := tryConnectingWithRetries(conString, logger, 3)
	if err != nil {
		logger.Errorw(constants.ErrorToStartDB, constants.Error, err.Error())
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw(constants.MsgToRetrieveSQLFromGorm, constants.Error, err.Error())
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		logger.Errorw(constants.FailedToPingDB, constants.Error, err.Error())
		return nil, err
	}

	return db, nil
}

func tryConnectingWithRetries(conString string, logger logger.Logger, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for tryConnect := 1; tryConnect <= maxRetries; tryConnect++ {
		logger.Infow(constants.MsgTryingStartsDB, constants.Try, tryConnect)
		db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		logger.Warnw(constants.ErrDBConnectionAttempt, constants.Error, err.Error())
		time.Sleep(3 * time.Second)
	}

	return nil, err
}

func Close(DB *gorm.DB, logger logger.Logger) {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Errorw(constants.MsgToRetrieveSQLFromGorm, constants.Error, err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Errorw(constants.ErrorToCloseDB, constants.Error, err.Error())
	} else {
		logger.Infow(constants.MsgPostgresConnectionClosed)
	}
}
