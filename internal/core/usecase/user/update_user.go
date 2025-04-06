package user

import (
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/constants"
	"time"
)

type UserUpdater interface {
	UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error)
	UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error)
}

func (s *UserService) UpdateUser(ctx domain.ContextControl, user domain.UserDomain) (domain.UserDomain, error) {
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

	return s.UserRepository.UpdateUser(ctx, user.ID, updateFields)
}

func (s *UserService) UpdateUserPassword(ctx domain.ContextControl, user domain.UserDomain, oldPassword, newPassword string) (domain.UserDomain, string, error) {
	userDB, err := s.UserRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToGetUserByID, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.SecurityHasher.ValidatePassword(userDB.Password, oldPassword); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	hashedPassword, err := s.SecurityHasher.HashPassword(newPassword)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToHashPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	fields := map[string]interface{}{
		constants.Password:  hashedPassword,
		constants.UpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.UserRepository.UpdateUser(ctx, user.ID, fields)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToUpdatePassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: user.ID}
	token, err := s.TokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}
	tokenDomain.Token = token

	if err := s.TokenService.Save(ctx, tokenDomain); err != nil {
		s.LoggerSugar.Errorw(constants.ErrorToSaveToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToSaveToken)
	}

	s.LoggerSugar.Infow(constants.SuccessPasswordUpdated, constants.UserID, updatedUser.ID)
	return updatedUser, tokenDomain.Token, nil
}
