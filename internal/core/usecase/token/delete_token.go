package token

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

//// Deleter defines an interface for deleting a token within a given context.
//// It accepts a TokenDomain object and returns an error if the operation fails.
// type Deleter interface {
//	Delete(ctx context.Context, token domain.TokenDomain) error
// }

// Delete removes the specified token from the repository and logs the result. Returns an error if the operation fails.
func (s *Service) Delete(ctx context.Context, token domain.TokenDomain) error {
	if err := s.tokenRepository.Delete(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, def.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessTokenDeleted, def.CtxUserID, token.UserID)

	return nil
}
