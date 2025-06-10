package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// Updater defines an interface for updating token information in a given context.
type Updater interface {
	Update(ctx context.Context, token domain.TokenDomain) error
}

// Update modifies the details of a given token in the repository and logs the operation's outcome. Returns an error if the update fails.
func (s *Service) Update(ctx context.Context, token domain.TokenDomain) error {
	if err := s.tokenRepository.Update(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToUpdateToken, constants.Error, err.Error())
		return err
	}
	s.logger.Infow(constants.SuccessTokenUpdated, constants.UserID, token.UserID)
	return nil
}
