package token

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Updater interface {
	Update(ctx context.Context, token domain.TokenDomain) error
}

func (s *TokenService) Update(ctx context.Context, token domain.TokenDomain) error {
	if err := s.tokenRepository.Update(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToUpdateToken, constants.Error, err.Error())
		return err
	}
	s.logger.Infow(constants.SuccessTokenUpdated, constants.UserID, token.UserID)
	return nil
}
