package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// Updater defines an interface for updating token information in a given context.
type Updater interface {
	Update(ctx context.Context, token entity.TokenDomain) error
}

// Update modifies the details of a given token in the repository and logs the operation's outcome. Returns an error if the update fails.
func (s *Service) Update(ctx context.Context, token entity.TokenDomain) error {
	if err := s.tokenRepository.Update(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToUpdateToken, def.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessTokenUpdated, def.CtxUserID, token.UserID)

	return nil
}
