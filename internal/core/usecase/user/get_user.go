package user

import (
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type UserRetriever interface {
	GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error)
	GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error)
	GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error)
	GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error)
}

func (s *UserService) GetUserByID(ctx domain.ContextControl, id uint64) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByID(ctx, id)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx domain.ContextControl, email string) (domain.UserDomain, error) {
	user, err := s.UserRepository.GetUserByEmail(ctx, email)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByEmail, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, user.ID)
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx domain.ContextControl, username string) (domain.UserDomain, error) {
	userDB, err := s.UserRepository.GetUserByUsername(ctx, username)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}
	s.LoggerSugar.Infow(constants.SuccessUserRetrieved, constants.UserID, userDB.ID)
	return userDB, nil
}

func (s *UserService) GetAllUsers(ctx domain.ContextControl) ([]domain.UserDomain, error) {
	users, err := s.UserRepository.GetAllUsers(ctx)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetAllUsers, constants.Error, err.Error())
		return nil, err
	}
	s.LoggerSugar.Infow(constants.SuccessUsersRetrieved, "count", len(users))
	return users, nil
}
