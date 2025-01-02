package service

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

const (
	ErrorToCreateUser = "error to create user"
	ErrorToGetUsers   = "error to get users"
	ErrorToGetUser    = "error to get user"
)

type UserService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	LoggerSugar                  *zap.SugaredLogger
}

func (service *UserService) CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.CreateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}
	return user, nil
}

func (service *UserService) GetAllUsers(contextControl domain.ContextControl) ([]domain.UserDomain, error) {

	users, err := service.UserDomainDataBaseRepository.GetAllUsers(contextControl)
	if err != nil {
		service.LoggerSugar.Errorw(ErrorToGetUsers, "error", err.Error())
		return []domain.UserDomain{}, err
	}
	return users, nil
}

func (service *UserService) GetUserByID(contextControl domain.ContextControl, ID uint64) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.GetUserByID(contextControl, ID)
	if err != nil {
		service.LoggerSugar.Errorw(ErrorToGetUser, "error", err.Error())
		return domain.UserDomain{}, err
	}
	return user, nil
}

func (service *UserService) UpdateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	user, err := service.UserDomainDataBaseRepository.UpdateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw(ErrorToCreateUser, "error", err.Error())
		return domain.UserDomain{}, err
	}
	return user, nil
}
