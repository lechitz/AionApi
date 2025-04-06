package user

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
)

type UserCreator interface {
	CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error)
}

func (s *UserService) CreateUser(ctx domain.ContextControl, user domain.UserDomain, password string) (domain.UserDomain, error) {
	user = s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToValidateCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToValidateCreateUser)
	}

	existingByUsername, err := s.UserRepository.GetUserByUsername(ctx, user.Username)
	if err == nil && existingByUsername.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.UsernameIsAlreadyInUse)
	}

	existingByEmail, err := s.UserRepository.GetUserByEmail(ctx, user.Email)
	if err == nil && existingByEmail.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.EmailIsAlreadyInUse)
	}

	hashPassword, err := s.SecurityHasher.HashPassword(password)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToHashPassword)
	}
	user.Password = hashPassword

	userDB, err := s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.LoggerSugar.Infow(constants.SuccessUserCreated, constants.UserID, userDB.ID)
	return userDB, nil
}
