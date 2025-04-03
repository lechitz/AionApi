package postgres

import (
	"fmt"
	"log"
	"time"

	"github.com/lechitz/AionApi/config"
	"github.com/lechitz/AionApi/internal/infrastructure/db/constants"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabaseConnection(config config.DBConfig, loggerSugar *zap.SugaredLogger) *gorm.DB {
	var err error
	conString := fmt.Sprintf(constants.MsgFormatConString, config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	loggerSugar.Infow(constants.MsgDBConnection, constants.Host, config.DBHost, constants.Port, config.DBPort, constants.DBName, config.DBName)

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
