package user

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type UserDeleter interface {
	SoftDeleteUser(ctx domain.ContextControl, id uint64) error
}

func (s *UserService) SoftDeleteUser(ctx domain.ContextControl, userID uint64) error {
	if err := s.UserRepository.SoftDeleteUser(ctx, userID); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSoftDeleteUser, constants.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{UserID: userID}
	if err := s.TokenService.Delete(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToDeleteToken, constants.Error, err.Error())
		return err
	}

	s.LoggerSugar.Infow(constants.SuccessUserSoftDeleted, constants.UserID, userID)
	return nil
}
