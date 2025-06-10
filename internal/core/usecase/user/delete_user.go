package user

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// Deleter defines a contract for deleting user records in a system.
// SoftDeleteUser deletes a user by marking them inactive or soft-deleted in the storage.
type Deleter interface {
	SoftDeleteUser(ctx context.Context, id uint64) error
}

// SoftDeleteUser performs a soft delete operation on a user identified by userID and deletes associated tokens. Returns an error if the operation fails.
func (s *Service) SoftDeleteUser(ctx context.Context, userID uint64) error {
	if err := s.userRepository.SoftDeleteUser(ctx, userID); err != nil {
		s.logger.Errorw(constants.ErrorToSoftDeleteUser, constants.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessUserSoftDeleted, constants.UserID, userID)
	return nil
}
