package service

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/pkg/utils"
	"github.com/lechitz/AionApi/ports/output"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type AuthService struct {
	UserDomainDataBaseRepository output.IUserDomainDataBaseRepository
	TokenStore                   output.ITokenStore
	LoggerSugar                  *zap.SugaredLogger
}

func (service *AuthService) Login(contextControl domain.ContextControl, userDomain domain.UserDomain) (domain.UserDomain, string, error) {
	userDB, err := service.UserDomainDataBaseRepository.GetUserByUsername(contextControl, userDomain)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	if err = utils.VerifyPassword(userDB.Password, userDomain.Password); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToVerifyPassword, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	token, err := service.TokenStore.CreateToken(contextControl.Context, userDB.ID)
	if err != nil {
		service.LoggerSugar.Errorw(constants.ErrorToCreateToken, "error", err.Error())
		return domain.UserDomain{}, "", err
	}

	return userDB, token, nil
}

func (service *AuthService) Logout(contextControl domain.ContextControl, userDomain domain.UserDomain, token string) error {
	if userDomain.ID == 0 {
		service.LoggerSugar.Errorw(constants.ErrorInvalidUserID, "userID", userDomain.ID)
		return fmt.Errorf(constants.ErrorInvalidUserID)
	}

	storedToken, err := service.TokenStore.GetTokenByUserID(contextControl.Context, userDomain.ID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			service.LoggerSugar.Warnw(constants.ErrorTokenNotFound, "userID", userDomain.ID)
			return fmt.Errorf(constants.ErrorTokenNotFound)
		}
		service.LoggerSugar.Errorw(constants.ErrorRetrieveTokenFromRedis, "error", err.Error(), "userID", userDomain.ID)
		return fmt.Errorf(constants.ErrorRetrieveTokenFromRedis)
	}

	if token != storedToken {
		service.LoggerSugar.Errorw(constants.ErrorTokenMismatch, "userID", userDomain.ID, "providedToken", token, "storedToken", storedToken)
		return fmt.Errorf(constants.ErrorTokenMismatch)
	}

	if err := service.TokenStore.DeleteTokenByUserID(contextControl.Context, userDomain.ID); err != nil {
		service.LoggerSugar.Errorw(constants.ErrorRevokeToken, "error", err.Error(), "userID", userDomain.ID)
		return err
	}

	service.LoggerSugar.Infow(constants.SuccessUserLoggedOut, "userID", userDomain.ID)
	return nil
}
