// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// CreateUser creates a new user with the given data and password, ensuring validations and unique constraints are met. Returns the created user or an error.
func (s *Service) CreateUser(ctx context.Context, user domain.UserDomain, password string) (domain.UserDomain, error) {
	s.normalizeUserData(&user)

	if err := s.validateCreateUserRequired(user, password); err != nil {
		s.logger.Errorw(constants.ErrorToValidateCreateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToValidateCreateUser)
	}

	existingByUsername, err := s.userStore.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw("DB error while checking username", commonkeys.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByUsername.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.UsernameIsAlreadyInUse)
	}

	existingByEmail, err := s.userStore.GetUserByEmail(ctx, user.Email)
	if err != nil {
		s.logger.Errorw("DB error while checking email", commonkeys.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToCreateUser)
	}

	if existingByEmail.ID != 0 {
		return domain.UserDomain{}, errors.New(constants.EmailIsAlreadyInUse)
	}

	hashedPassword, err := s.hashStore.HashPassword(password)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, commonkeys.Error, err.Error())
		return domain.UserDomain{}, errors.New(constants.ErrorToHashPassword)
	}

	user.Password = hashedPassword

	userDB, err := s.userStore.CreateUser(ctx, user)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserCreated, commonkeys.UserID, strconv.FormatUint(userDB.ID, 10))

	return userDB, nil
}
