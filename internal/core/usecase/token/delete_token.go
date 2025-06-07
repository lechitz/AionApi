package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

type Deleter interface {
	Delete(ctx context.Context, token domain.TokenDomain) error
}

func (s *TokenService) Delete(ctx context.Context, token domain.TokenDomain) error {
	if err := s.tokenRepository.Delete(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}
	s.logger.Infow(constants.SuccessTokenDeleted, constants.UserID, token.UserID)
	return nil
}
