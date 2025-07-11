package token

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/adapters/secondary/security"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// CreateToken generates a new token for the provided user, saves it in the repository, and returns the signed token or an error.
func (s *Service) CreateToken(ctx context.Context, tokenDomain domain.TokenDomain) (string, error) {
	value, err := s.tokenRepository.Get(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetToken, commonkeys.Error, err.Error())
		return "", err
	}

	if value != "" {
		if err := s.tokenRepository.Delete(ctx, tokenDomain); err != nil {
			s.logger.Errorw(constants.ErrorToDeleteToken, commonkeys.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := security.GenerateToken(tokenDomain.UserID, s.secretKey)
	if err != nil {
		s.logger.Errorw(constants.ErrorToAssignToken, commonkeys.Error, err.Error())
		return "", err
	}

	tokenDomain.Token = signedToken

	err = s.tokenRepository.Save(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, commonkeys.Error, err.Error())
		return "", err
	}

	s.logger.Infow(constants.SuccessTokenCreated, commonkeys.UserID, strconv.FormatUint(tokenDomain.UserID, 10))
	return signedToken, nil
}
