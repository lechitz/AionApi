package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/infrastructure/security"
)

type Creator interface {
	CreateToken(ctx domain.ContextControl, token domain.TokenDomain) (string, error)
}

func (s *TokenService) CreateToken(ctx domain.ContextControl, tokenDomain domain.TokenDomain) (string, error) {
	if _, err := s.TokenRepository.Get(ctx, tokenDomain); err == nil {
		if err := s.TokenRepository.Delete(ctx, tokenDomain); err != nil {
			s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := security.GenerateToken(tokenDomain.UserID, s.ConfigToken.SecretKey)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToAssignToken, constants.Error, err.Error())
		return "", err
	}

	tokenDomain.Token = signedToken

	if err := s.TokenRepository.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return "", err
	}

	s.LoggerSugar.Infow(constants.SuccessTokenCreated, constants.UserID, tokenDomain.UserID)
	return signedToken, nil
}
