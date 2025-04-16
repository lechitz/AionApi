package token

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"github.com/lechitz/AionApi/internal/infrastructure/security"
)

type Creator interface {
	CreateToken(ctx context.Context, token domain.TokenDomain) (string, error)
}

func (s *TokenService) CreateToken(ctx context.Context, tokenDomain domain.TokenDomain) (string, error) {
	if _, err := s.tokenRepository.Get(ctx, tokenDomain); err == nil {
		if err := s.tokenRepository.Delete(ctx, tokenDomain); err != nil {
			s.logger.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := security.GenerateToken(tokenDomain.UserID, s.configToken.SecretKey)
	if err != nil {
		s.logger.Errorw(constants.ErrorToAssignToken, constants.Error, err.Error())
		return "", err
	}

	tokenDomain.Token = signedToken

	if err := s.tokenRepository.Save(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return "", err
	}

	s.logger.Infow(constants.SuccessTokenCreated, constants.UserID, tokenDomain.UserID)
	return signedToken, nil
}
