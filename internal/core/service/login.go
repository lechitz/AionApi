package service

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/ports/output"
	"go.uber.org/zap"
)

const (
	ErrorToLogin = "error to login"
)

type LoginService struct {
	LoginDomainDataBaseRepository output.ILoginDomainDataBaseRepository
	LoggerSugar                   *zap.SugaredLogger
}

func (service *LoginService) GetUserByUsername(contextControl domain.ContextControl, user domain.LoginDomain) (domain.LoginDomain, error) {

	login, err := service.LoginDomainDataBaseRepository.GetUserByUsername(contextControl, user)
	if err != nil {
		service.LoggerSugar.Errorw(ErrorToLogin, "error", err.Error())
		return domain.LoginDomain{}, err
	}
	return login, nil
}
