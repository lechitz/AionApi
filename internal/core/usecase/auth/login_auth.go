// Package auth contains use cases for authenticating users and generating tokens.
package auth

import (
	"context"
	"errors"
	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
// Returns the user data, token, and error if any occurs.
func (s *Service) Login(ctx context.Context, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error) {
	userDB, err := s.userRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if userDB.ID == 0 {
		s.logger.Warnw(constants.UserNotFoundOrInvalidCredentials, def.CtxUsername, user.Username)
		return domain.UserDomain{}, "", errors.New(constants.UserNotFoundOrInvalidCredentials)
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, passwordReq); err != nil {
		s.logger.Warnw(constants.ErrorToCompareHashAndPassword, def.CtxUsername, user.Username)
		return domain.UserDomain{}, "", errors.New(constants.InvalidCredentials)
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	newToken, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, def.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.logger.Infow(constants.SuccessToLogin, def.CtxUserID, userDB.ID, def.CtxToken, newToken)

	return userDB, newToken, nil
}
