// Package auth contains use cases for authenticating users and generating tokens.
package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
func (s *Service) Login(ctx context.Context, username, password string) (domain.User, string, error) {
	userDomain, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, commonkeys.Error, err.Error())
		return domain.User{}, "", errors.New(constants.ErrorToGetUserByUserName)
	}

	if userDomain.ID == 0 {
		s.logger.Warnw(constants.UserNotFoundOrInvalidCredentials, commonkeys.Username, userDomain.Username)
		return domain.User{}, "", errors.New(constants.UserNotFoundOrInvalidCredentials)
	}

	if err := s.securityHasher.Compare(userDomain.Password, password); err != nil {
		s.logger.Warnw(constants.ErrorToCompareHashAndPassword, commonkeys.Username, userDomain.Username)
		return domain.User{}, "", errors.New(constants.InvalidCredentials)
	}

	tokenDomain, err := s.tokenService.CreateToken(ctx, userDomain.ID)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, commonkeys.Error, err.Error())
		return domain.User{}, "", errors.New(constants.ErrorToCreateToken)
	}

	s.logger.Infow(constants.SuccessToLogin, commonkeys.UserID, strconv.FormatUint(userDomain.ID, 10))

	return userDomain, tokenDomain.Token, nil
}
