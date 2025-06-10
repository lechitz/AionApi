package user

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type UserCreator interface {
	CreateUser(
		ctx context.Context,
		user domain.UserDomain,
		password string,
	) (domain.UserDomain, error)
}

func (s *UserService) CreateUser(
	ctx context.Context,
	user domain.UserDomain,
	password string,
) (domain.UserDomain, error) {
	user = s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.logger.Errorw(constants.ErrorToValidateCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToValidateCreateUser)
	}

	if existingByUsername, err := s.userRepository.GetUserByUsername(ctx, user.Username); err == nil &&
		existingByUsername.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.UsernameIsAlreadyInUse)
	}

	if existingByEmail, err := s.userRepository.GetUserByEmail(ctx, user.Email); err == nil &&
		existingByEmail.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.EmailIsAlreadyInUse)
	}

	hashedPassword, err := s.securityHasher.HashPassword(password)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToHashPassword)
	}

	user.Password = hashedPassword

	userDB, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserCreated, constants.UserID, userDB.ID)
	return userDB, nil
}
