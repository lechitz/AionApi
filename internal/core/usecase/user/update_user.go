// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// Updater defines methods to update user information and change user passwords in the system.
// UpdateUser updates a user's details in the system and returns the updated user or an error.
// UpdateUserPassword changes a user's password, verifying the old password, and returns the updated user, a confirmation, or an error.
type Updater interface {
	UpdateUser(ctx context.Context, user entity.UserDomain) (entity.UserDomain, error)
	UpdateUserPassword(ctx context.Context, user entity.UserDomain, oldPassword, newPassword string) (entity.UserDomain, string, error)
}

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, user entity.UserDomain) (entity.UserDomain, error) {
	updateFields := make(map[string]interface{})

	if user.Name != "" {
		updateFields[constants.Name] = user.Name
	}

	if user.Username != "" {
		updateFields[constants.Username] = user.Username
	}

	if user.Email != "" {
		updateFields[constants.Email] = user.Email
	}

	if len(updateFields) == 0 {
		return entity.UserDomain{}, errors.New(constants.ErrorNoFieldsToUpdate)
	}

	updateFields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdateUser, def.Error, err.Error())
		return entity.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserUpdated, def.CtxUserID, updatedUser.ID)

	return updatedUser, nil
}

// UpdateUserPassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdateUserPassword(ctx context.Context, user entity.UserDomain, oldPassword, newPassword string) (entity.UserDomain, string, error) {
	userDB, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, oldPassword); err != nil {
		s.logger.Errorw(constants.ErrorToCompareHashAndPassword, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	hashedPassword, err := s.securityHasher.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		constants.Password:  hashedPassword,
		constants.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdatePassword, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	tokenDomain := entity.TokenDomain{UserID: user.ID}

	token, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	tokenDomain.Token = token

	if err := s.tokenService.Save(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, def.Error, err.Error())
		return entity.UserDomain{}, "", errors.New(constants.ErrorToSaveToken)
	}

	s.logger.Infow(constants.SuccessPasswordUpdated, def.CtxUserID, updatedUser.ID)

	return updatedUser, token, nil
}
