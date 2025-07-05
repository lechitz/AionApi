package user

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// SoftDeleteUser performs a soft delete operation on a user identified by userID and deletes associated tokens. Returns an error if the operation fails.
func (s *Service) SoftDeleteUser(ctx context.Context, userID uint64) error {
	if err := s.userRepository.SoftDeleteUser(ctx, userID); err != nil {
		s.logger.Errorw(constants.ErrorToSoftDeleteUser, def.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, def.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessUserSoftDeleted, def.CtxUserID, userID)
	return nil
}
