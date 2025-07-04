// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// Creator defines the interface for creating a new user in the system. It returns the created user or an error.
type Creator interface {
	CreateUser(ctx context.Context, user entity.UserDomain, password string) (entity.UserDomain, error)
}

// CreateUser creates a new user with the given data and password, ensuring validations and unique constraints are met. Returns the created user or an error.
func (s *Service) CreateUser(ctx context.Context, user entity.UserDomain, password string) (entity.UserDomain, error) {
	user = s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.logger.Errorw(constants.ErrorToValidateCreateUser, def.Error, err.Error())
		return entity.UserDomain{}, errors.New(constants.ErrorToValidateCreateUser)
	}

	existingByUsername, err := s.userRepository.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw("DB error while checking username", def.Error, err)
		return entity.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByUsername.ID != 0 {
		return entity.UserDomain{}, errors.New(constants.UsernameIsAlreadyInUse)
	}

	existingByEmail, err := s.userRepository.GetUserByEmail(ctx, user.Email)
	if err != nil {
		s.logger.Errorw("DB error while checking email", def.Error, err)
		return entity.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByEmail.ID != 0 {
		return entity.UserDomain{}, errors.New(constants.EmailIsAlreadyInUse)
	}

	hashedPassword, err := s.securityHasher.HashPassword(password)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, def.Error, err.Error())
		return entity.UserDomain{}, errors.New(constants.ErrorToHashPassword)
	}

	user.Password = hashedPassword

	userDB, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateUser, def.Error, err.Error())
		return entity.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserCreated, def.CtxUserID, userDB.ID)

	return userDB, nil
}
