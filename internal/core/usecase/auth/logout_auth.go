package auth

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type SessionRevoker interface {
	Logout(ctx domain.ContextControl, token string) error
}

func (s *AuthService) Logout(ctx domain.ContextControl, token string) error {
	userID, _, err := s.TokenService.Check(ctx, token)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCheckToken, constants.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToRevokeToken, constants.Error, err.Error(), constants.UserID, userID)
		return err
	}

	s.LoggerSugar.Infow(constants.SuccessUserLoggedOut, constants.UserID, userID)
	return nil
}
