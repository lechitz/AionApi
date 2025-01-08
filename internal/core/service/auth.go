package service

import (
	"github.com/lechitz/AionApi/adapters/input/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

type AuthService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	LoggerSugar                  *zap.SugaredLogger
}

func (service *AuthService) Login(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, string, error) {

	userDB, err := service.UserDomainDataBaseRepository.GetUserByUsername(contextControl, userDomain)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByUsername, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	if err = utils.VerifyPassword(userDB.Password, userDomain.Password); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToVerifyPassword, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	token, err := utils.CreateToken(userDB.ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToCreateToken, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	return userDB, token, nil
}
