package auth

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Authenticator interface {
	Login(ctx domain.ContextControl, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error)
}

func (s *AuthService) Login(ctx domain.ContextControl, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {
	userDB, err := s.UserRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.PasswordHasher.ComparePasswords(userDB.Password, passwordReq); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	newToken, err := s.TokenService.Create(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain.Token = newToken

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.LoggerSugar.Infow(constants.SuccessToLogin, constants.UserID, userDB.ID)
	return userDB, tokenDomain.Token, nil
}
