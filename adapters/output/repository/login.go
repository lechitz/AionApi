package repository

import (
	"github.com/jinzhu/copier"
	"github.com/lechitz/AionApi/internal/core/domain"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoginPostgresDB struct {
	DB          *gorm.DB
	LoggerSugar *zap.SugaredLogger
}

type LoginDB struct {
	ID       uint64 `gorm:"primaryKey, column:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func NewLoginPostgresDB(gormDB *gorm.DB, loggerSugar *zap.SugaredLogger) LoginPostgresDB {
	return LoginPostgresDB{
		DB:          gormDB,
		LoggerSugar: loggerSugar,
	}
}

func (LoginDB) TableName() string {
	return "aion_api.users"
}

func (lg LoginDB) CopyToLoginDomain() domain.LoginDomain {
	return domain.LoginDomain{
		ID:       lg.ID,
		Username: lg.Username,
		Password: lg.Password,
	}
}

func (lg LoginPostgresDB) GetUserByUsername(contextControl domain.ContextControl, loginDomain domain.LoginDomain) (domain.LoginDomain, error) {

	var loginDB LoginDB
	copier.Copy(&loginDB, &loginDomain)

	if err := lg.DB.WithContext(contextControl.Context).
		Select("id, username, password").
		Where("username = ?", loginDB.Username).
		First(&loginDB).Error; err != nil {
		return domain.LoginDomain{}, err
	}

	return loginDB.CopyToLoginDomain(), nil
}
