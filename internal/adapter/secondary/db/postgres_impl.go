// Package db provides database connection and management functions.
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/lechitz/AionApi/internal/platform/config"
	"github.com/lechitz/AionApi/internal/platform/ports/output/logger"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MsgFormatConString is the template for the database connection string.
const MsgFormatConString = "host=%s ports=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s"

// MsgTryingStartsDB indicates an attempt to establish a database connection.
const MsgTryingStartsDB = "trying to establish Database connection"

// ErrorToStartDB indicates a failure to start the PostgresSQL database.
const ErrorToStartDB = "error to start the postgres db"

// MsgToRetrieveSQLFromGorm indicates a failure to retrieve the SQL DB from Gorm.
const MsgToRetrieveSQLFromGorm = "failed to retrieve SQL DB from Gorm"

// ErrorToCloseDB indicates a failure to close the PostgresSQL database.
const ErrorToCloseDB = "error to close the postgres db"

// MsgDBConnection indicates the initialization of a database connection.
const MsgDBConnection = "initializing Database connection"

// FailedToPingDB indicates a failed ping to the database.
const FailedToPingDB = "failed to ping DB"

// ErrDBConnectionAttempt indicates a failed database connection attempt.
const ErrDBConnectionAttempt = "connection attempt failed"

// MsgPostgresConnectionClosed indicates that the PostgresSQL connection was closed successfully.
const MsgPostgresConnectionClosed = "Database connection closed successfully"

// NewConnection initializes a database connection using the provided configuration and contextlogger. Returns a Gorm DB instance or an error.
func NewConnection(appCtx context.Context, cfg config.DBConfig, logger logger.ContextLogger) (*gorm.DB, error) {
	db, err := tryConnectingWithRetries(cfg, logger)
	if err != nil {
		logger.Errorw(ErrorToStartDB, commonkeys.Error, err.Error())
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw(MsgToRetrieveSQLFromGorm, commonkeys.Error, err.Error())
		return nil, err
	}

	if err := sqlDB.PingContext(appCtx); err != nil {
		logger.Errorw(FailedToPingDB, commonkeys.Error, err.Error())
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
func tryConnectingWithRetries(cfg config.DBConfig, logger logger.ContextLogger) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	conString := fmt.Sprintf(MsgFormatConString, cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.TimeZone)
	logger.Infow(MsgDBConnection, commonkeys.DBHost, cfg.Host, commonkeys.DBPort, cfg.Port, commonkeys.DBName, cfg.Name)

	for tryConnect := 1; tryConnect <= cfg.MaxRetries; tryConnect++ {
		logger.Infow(MsgTryingStartsDB, commonkeys.DBTryConnectingWithRetries, tryConnect)
		db, err = gorm.Open(postgres.Open(conString), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		logger.Warnw(ErrDBConnectionAttempt, commonkeys.Error, err.Error())
		time.Sleep(cfg.RetryInterval)
	}

	return nil, err
}

// Close terminates the database connection and logs success or error messages using the provided contextlogger.
func Close(db *gorm.DB, logger logger.ContextLogger) {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorw(MsgToRetrieveSQLFromGorm, commonkeys.Error, err.Error())
		return
	}

	if err := sqlDB.Close(); err != nil {
		logger.Errorw(ErrorToCloseDB, commonkeys.Error, err.Error())
	} else {
		logger.Infow(MsgPostgresConnectionClosed)
	}
}
