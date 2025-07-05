package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/adapters/secondary/security"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// CreateToken generates a new token for the provided user, saves it in the repository, and returns the signed token or an error.
func (s *Service) CreateToken(ctx context.Context, tokenDomain domain.TokenDomain) (string, error) {
	if _, err := s.tokenRepository.Get(ctx, tokenDomain); err == nil {
		if err := s.tokenRepository.Delete(ctx, tokenDomain); err != nil {
			s.logger.Errorw(constants.ErrorToDeleteToken, def.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := security.GenerateToken(tokenDomain.UserID, s.configToken.SecretKey)
	if err != nil {
		s.logger.Errorw(constants.ErrorToAssignToken, def.Error, err.Error())
		return "", err
	}

	tokenDomain.Token = signedToken

	if err := s.tokenRepository.Save(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, def.Error, err.Error())
		return "", err
	}

	s.logger.Infow(constants.SuccessTokenCreated, def.CtxUserID, tokenDomain.UserID)
	return signedToken, nil
}
