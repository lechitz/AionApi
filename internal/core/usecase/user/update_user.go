package user

import (
	"context"
	"errors"
	"time"

	"github.com/lechitz/AionApi/internal/core/usecase/user/constants"

	"github.com/lechitz/AionApi/internal/core/domain"
)

type UserUpdater interface {
	UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
}

func (s *UserService) UpdateUser(ctx context.Context, user domain.UserDomain) (domain.UserDomain, error) {
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
		s.logger.Errorw(constants.ErrorToUpdateUser, constants.Error, err.Error())
		return domain.UserDomain{}, err
	}

	s.logger.Infow(constants.SuccessUserUpdated, constants.UserID, updatedUser.ID)
	return updatedUser, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, oldPassword); err != nil {
		s.logger.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := s.securityHasher.HashPassword(newPassword)
	if err != nil {
		s.logger.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		constants.Password:  hashedPassword,
		constants.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.logger.Errorw(constants.ErrorToUpdatePassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}
	token, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := s.tokenService.Save(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToSaveToken)
	}

	s.logger.Infow(constants.SuccessPasswordUpdated, constants.UserID, updatedUser.ID)
	return updatedUser, token, nil
}
