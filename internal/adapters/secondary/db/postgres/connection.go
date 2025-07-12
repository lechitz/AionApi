// Package db provides database connection and management functions.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres/constants"

	"github.com/lechitz/AionApi/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewConnection initializes a database connection using the provided configuration and contextlogger. Returns a Gorm DB instance or an error.
func NewConnection(appCtx context.Context, cfg config.DBConfig, logger output.ContextLogger) (*gorm.DB, error) {
	db, err := tryConnectingWithRetries(cfg, logger)
	if err != nil {
		logger.Errorw(constants.ErrorToStartDB, commonkeys.Error, err.Error())
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw(constants.MsgToRetrieveSQLFromGorm, commonkeys.Error, err.Error())
		return nil, err
	}

	if err := sqlDB.PingContext(appCtx); err != nil {
		logger.Errorw(constants.FailedToPingDB, commonkeys.Error, err.Error())
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	return db, nil
}

// tryConnectingWithRetries attempts to establish a database connection with retries.
func tryConnectingWithRetries(cfg config.DBConfig, logger output.ContextLogger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	conString := fmt.Sprintf(constants.MsgFormatConString, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.TimeZone)
	logger.Infow(constants.MsgDBConnection, commonkeys.DBHost, cfg.Host, commonkeys.DBPort, cfg.Port, commonkeys.DBName, cfg.Name)

	for tryConnect := 1; tryConnect <= cfg.MaxRetries; tryConnect++ {
		logger.Infow(constants.MsgTryingStartsDB, commonkeys.DBTryConnectingWithRetries, tryConnect)
		db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		logger.Warnw(constants.ErrDBConnectionAttempt, commonkeys.Error, err.Error())
		time.Sleep(cfg.RetryInterval)
	}

	return nil, err
}

// Close terminates the database connection and logs success or error messages using the provided contextlogger.
func Close(db *gorm.DB, logger output.ContextLogger) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw(constants.MsgToRetrieveSQLFromGorm, commonkeys.Error, err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Errorw(constants.ErrorToCloseDB, commonkeys.Error, err.Error())
	} else {
		logger.Infow(constants.MsgPostgresConnectionClosed)
	}
}
