package token

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/common"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
)

// Delete removes the specified token from the repository and logs the result. Returns an error if the operation fails.
func (s *Service) Delete(ctx context.Context, token domain.TokenDomain) error {
	if err := s.tokenRepository.Delete(ctx, token); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, common.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessTokenDeleted, common.UserID, strconv.FormatUint(token.UserID, 10))

	return nil
}
