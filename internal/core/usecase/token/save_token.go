package token

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type Saver interface {
	Save(ctx domain.ContextControl, token domain.TokenDomain) error
}

func (s *TokenService) Save(ctx domain.ContextControl, token domain.TokenDomain) error {
	if err := s.TokenRepository.Save(ctx, token); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return err
	}
	s.LoggerSugar.Infow(constants.SuccessTokenCreated, constants.UserID, token.UserID)
	return nil
}
