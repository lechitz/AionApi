// Package user contains use cases for managing users in the system.
package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"
)

// UpdateUser updates an existing user's attributes based on the provided data. Returns the updated user or an error if the operation fails.
func (s *Service) UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error) {
	updateFields := make(map[string]interface{})

	if user.Name != "" {
		updateFields[commonkeys.Name] = user.Name
	}

	if user.Username != "" {
		updateFields[commonkeys.Username] = user.Username
	}

	if user.Email != "" {
		updateFields[commonkeys.Email] = user.Email
	}

	if len(updateFields) == 0 {
		return domain.UserDomain{}, errors.New(constants.ErrorNoFieldsToUpdate)
	}

	updateFields[constants.UpdatedAt] = time.Now().UTC()

	updatedUser, err := s.userStore.UpdateUser(ctx, user.ID, updateFields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdateUser, commonkeys.Error, err.Error())
		return domain.UserDomain{}, fmt.Errorf("%s: %w", constants.ErrorToUpdateUser, err)
	}

	s.logger.Infow(constants.SuccessUserUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10), commonkeys.UserUpdatedFields, updateFields)

	return updatedUser, nil
}

// UpdateUserPassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.userStore.GetUserByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))

		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToGetUserByID, err)
	}

	if err := s.hashStore.ValidatePassword(userDB.Password, oldPassword); err != nil {
		s.logger.Errorw(constants.ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))

		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hashStore.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))

		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToHashPassword, err)
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userStore.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))

		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToUpdatePassword, err)
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}

	token, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))

		return domain.UserDomain{}, "", fmt.Errorf("%s: %w", constants.ErrorToCreateToken, err)
	}

	s.logger.Infow(constants.SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return updatedUser, token, nil
}
