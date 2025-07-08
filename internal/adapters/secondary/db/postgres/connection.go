// Package db provides database connection and management functions.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/core/ports/output"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/postgres/constants"

	"github.com/lechitz/AionApi/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewDatabaseConnection initializes a database connection using the provided configuration and logger. Returns a Gorm DB instance or an error.
func NewDatabaseConnection(appCtx context.Context, cfg config.DBConfig, logger output.Logger) (*gorm.DB, error) {
	conString := fmt.Sprintf(constants.MsgFormatConString, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	// TODO: avaliar variáveis abaixo.

	logger.Infow(constants.MsgDBConnection, constants.Host, cfg.Host, commonkeys.Port, cfg.Port, constants.DBName, cfg.Name)

	db, err := tryConnectingWithRetries(conString, logger, 3)
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
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute) // TODO: avaliar se vale apena passar pra time.Duration.
	}

	return db, nil
}

// tryConnectingWithRetries attempts to establish a database connection with retries.
func tryConnectingWithRetries(conString string, logger output.Logger, maxRetries int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	for tryConnect := 1; tryConnect <= maxRetries; tryConnect++ {
		logger.Infow(constants.MsgTryingStartsDB, constants.Try, tryConnect) // TODO: avaliar se vale apena passar pra o Try para commonkeys.
		db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		logger.Warnw(constants.ErrDBConnectionAttempt, commonkeys.Error, err.Error())
		time.Sleep(3 * time.Second) // TODO: avaliar se vale apena ir para variáveis de ambiente.
	}

	return nil, err
}

// Close terminates the database connection and logs success or error messages using the provided logger.
func Close(db *gorm.DB, logger output.Logger) {
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
