// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"time"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error) {
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
		return domain.UserDomain{}, errors.New(constants.ErrorNoFieldsToUpdate)
	}

	updateFields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdateUser, def.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserUpdated, def.CtxUserID, updatedUser.ID)

	return updatedUser, nil
}

// UpdateUserPassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, oldPassword); err != nil {
		s.logger.Errorw(constants.ErrorToCompareHashAndPassword, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := s.securityHasher.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		constants.Password:  hashedPassword,
		constants.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdatePassword, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}

	token, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain.Token = token

	if err := s.tokenService.Save(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, def.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToSaveToken)
	}

	s.logger.Infow(constants.SuccessPasswordUpdated, def.CtxUserID, updatedUser.ID)

	return updatedUser, token, nil
}
