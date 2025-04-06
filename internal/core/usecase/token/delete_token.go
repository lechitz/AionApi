package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Deleter interface {
	Delete(ctx domain.ContextControl, token domain.TokenDomain) error
}

func (s *TokenService) Delete(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Delete(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(constants.SuccessTokenDeleted, constants.UserID, token.UserID)
	return nil
}
