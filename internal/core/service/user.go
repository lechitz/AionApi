package service

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

type UserService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	LoggerSugar                  *zap.SugaredLogger
}

func (service *UserService) CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.CreateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw("Error to create user", "error", err.Error())
		return domain.UserDomain{}, err
	}
	return user, nil
}
