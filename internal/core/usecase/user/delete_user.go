package user

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/common"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// SoftDeleteUser performs a soft delete operation on a user identified by userID and deletes associated tokens. Returns an error if the operation fails.
func (s *Service) SoftDeleteUser(ctx context.Context, userID uint64) error {
	if err := s.userStore.SoftDeleteUser(ctx, userID); err != nil {
		s.logger.Errorw(constants.ErrorToSoftDeleteUser, common.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToDeleteToken, common.Error, err.Error())
		return err
	}

	s.logger.Infow(constants.SuccessUserSoftDeleted, common.UserID, strconv.FormatUint(userID, 10))
	return nil
}
