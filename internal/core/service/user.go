package service

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

type UserService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	LoggerSugar                  *zap.SugaredLogger
}

func (service *UserService) CreateUser(contextControl domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {

	if err := validateUser(user); err != nil {
		service.LoggerSugar.Errorw("Error to validate user", "error", err.Error())
		return domain.UserDomain{}, err
	}

	user, err := service.UserDomainDataBaseRepository.CreateUser(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw("Error to create user", "error", err.Error())
		return domain.UserDomain{}, err
	}
	return user, nil
}

func validateUser(user domain.UserDomain) error {
	if user.Name == "" {
		return errors.New("invalid name")
	}
	if user.Username == "" {
		return errors.New("invalid username")
	}
	if user.Email == "" {
		return errors.New("invalid email")
	}
	if user.Password == "" {
		return errors.New("invalid password")
	}
	return nil
}
