package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Creator interface {
	CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
}

func (s *TokenService) CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error) {
	if _, err := s.TokenRepository.Get(ctx, token); err == nil {
		if err := s.TokenRepository.Delete(ctx, token); err != nil {
			s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := generateToken(token.UserID)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToAssignToken, constants.Error, err.Error())
		return "", err
	}

	token.Token = signedToken

	if err := s.TokenRepository.Save(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return "", err
	}

	s.LoggerSugar.Infow(constants.SuccessTokenCreated, constants.UserID, token.UserID)
	return signedToken, nil
}
