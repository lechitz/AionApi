package user

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

type UserDeleter interface {
	SoftDeleteUser(ctx context.Context, id uint64) error
}

func (s *UserService) SoftDeleteUser(ctx context.Context, userID uint64) error {
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
