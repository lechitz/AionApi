package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Updater interface {
	Update(ctx domain.ContextControl, token domain.TokenDomain) error
}

func (s *TokenService) Update(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.tokenRepository.Update(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToUpdateToken, constants.Error, err.Error())
		return err
	}
	s.logger.Infow(constants.SuccessTokenUpdated, constants.UserID, token.UserID)
	return nil
}
