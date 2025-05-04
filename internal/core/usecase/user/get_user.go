package user

import (
	"context"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

type UserRetriever interface {
	GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error)
	GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error)
	GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error)
	GetAllUsers(ctx context.Context) ([]domain.UserDomain, error)
}

func (s *UserService) GetUserByID(ctx context.Context, userID uint64) (domain.UserDomain, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.logger.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (domain.UserDomain, error) {
	user, err := s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByEmail, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.logger.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (domain.UserDomain, error) {
	userDB, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.logger.Infow(constants.SuccessUserRetrieved, constants.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]domain.UserDomain, error) {
	users, err := s.userRepository.GetAllUsers(ctx)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetAllUsers, constants.Error, err.Error())
		return nil, err
	}
	s.logger.Infow(constants.SuccessUsersRetrieved, constants.Users, len(users))
	return users, nil
}
