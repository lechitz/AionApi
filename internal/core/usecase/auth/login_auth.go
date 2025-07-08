// Package auth contains use cases for authenticating users and generating tokens.
package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
func (s *Service) Login(ctx context.Context, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {
	userDB, err := s.userRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, commonkeys.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToGetUserByUserName)
	}

	if userDB.ID == 0 {
		s.logger.Warnw(constants.UserNotFoundOrInvalidCredentials, commonkeys.Username, user.Username)
		return domain.UserDomain{}, "", errors.New(constants.UserNotFoundOrInvalidCredentials)
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, passwordReq); err != nil {
		s.logger.Warnw(constants.ErrorToCompareHashAndPassword, commonkeys.Username, user.Username)
		return domain.UserDomain{}, "", errors.New(constants.InvalidCredentials)
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	newToken, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, commonkeys.Error, err.Error())
		return domain.UserDomain{}, "", errors.New(constants.ErrorToCreateToken)
	}

	s.logger.Infow(constants.SuccessToLogin, commonkeys.UserID, strconv.FormatUint(userDB.ID, 10), commonkeys.Token, newToken)

	return userDB, newToken, nil
}
