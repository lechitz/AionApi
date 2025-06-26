package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// Saver defines an interface for persisting token information in a given storage.
// It expects a context and a TokenDomain object to perform the save operation.
// Save method returns an error if the persistence operation fails.
type Saver interface {
	Save(ctx context.Context, token entity.TokenDomain) error
}

// Save persists the provided token in the repository and logs success or error messages. Returns an error if saving fails.
func (s *Service) Save(ctx context.Context, token entity.TokenDomain) error {
	if err := s.tokenRepository.Save(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, def.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessTokenCreated, def.CtxUserID, token.UserID)

	return nil
}
