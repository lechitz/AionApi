// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// CreateUser creates a new user with the given data and password, ensuring validations and unique constraints are met. Returns the created user or an error.
func (s *Service) CreateUser(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, error) {
	user = s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.logger.Errorw(constants.ErrorToValidateCreateUser, def.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToValidateCreateUser)
	}

	existingByUsername, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw("DB error while checking username", def.Error, err)
		return domain.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByUsername.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.UsernameIsAlreadyInUse)
	}

	existingByEmail, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		s.logger.Errorw("DB error while checking email", def.Error, err)
		return domain.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByEmail.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.EmailIsAlreadyInUse)
	}

	hashedPassword, err := s.securityHasher.HashPassword(password)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, def.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToHashPassword)
	}

	user.Password = hashedPassword

	userDB, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateUser, def.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserCreated, def.CtxUserID, userDB.ID)

	return userDB, nil
}
