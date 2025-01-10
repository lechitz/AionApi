package service

import (
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/adapters/input/constants"
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
		service.LoggerSugar.Errorw(constants.ErrorToGetUserByUsername, "error", err.Error())
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
		service.LoggerSugar.Errorw("Invalid userID provided for logout", "userID", userDomain.ID)
		return fmt.Errorf("invalid userID")
	}

	storedToken, err := service.TokenStore.GetTokenByUserID(contextControl.Context, userDomain.ID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			service.LoggerSugar.Warnw("No token found in Redis for user during logout", "userID", userDomain.ID)
			return fmt.Errorf("token not found")
		}
		service.LoggerSugar.Errorw("Error retrieving token from Redis", "error", err.Error(), "userID", userDomain.ID)
		return fmt.Errorf("error retrieving token from Redis")
	}

	if token != storedToken {
		service.LoggerSugar.Errorw("Token mismatch during logout", "userID", userDomain.ID, "providedToken", token, "storedToken", storedToken)
		return fmt.Errorf("provided token does not match stored token")
	}

	if err := service.TokenStore.DeleteTokenByUserID(contextControl.Context, userDomain.ID); err != nil {
		service.LoggerSugar.Errorw("Failed to revoke token during logout", "error", err.Error(), "userID", userDomain.ID)
		return err
	}

	service.LoggerSugar.Infow("User logged out successfully", "userID", userDomain.ID)
	return nil
}
