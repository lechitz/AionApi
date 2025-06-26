package auth

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/core/domain/entity"
	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Authenticator defines an interface for handling user authentication operations.
// Login validates user credentials and returns user details with a session token.
type Authenticator interface {
	Login(ctx context.Context, user entity.UserDomain, passwordReq string) (entity.UserDomain, string, error)
}

// Login authenticates a user by validating credentials and generates a new token if valid.
// Returns the user data, token, and error if any occurs.
func (s *Service) Login(ctx context.Context, user entity.UserDomain, passwordReq string) (entity.UserDomain, string, error) {
	userDB, err := s.userRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	if userDB.ID == 0 {
		s.logger.Warnw(constants.UserNotFoundOrInvalidCredentials, def.CtxUsername, user.Username)
		return entity.UserDomain{}, "", errors.New(constants.UserNotFoundOrInvalidCredentials)
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, passwordReq); err != nil {
		s.logger.Warnw(constants.ErrorToCompareHashAndPassword, def.CtxUsername, user.Username)
		return entity.UserDomain{}, "", errors.New(constants.InvalidCredentials)
	}

	tokenDomain := entity.TokenDomain{UserID: userDB.ID}

	newToken, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, def.Error, err.Error())
		return entity.UserDomain{}, "", err
	}

	s.logger.Infow(constants.SuccessToLogin, def.CtxUserID, userDB.ID, def.CtxToken, newToken)

	return userDB, newToken, nil
}
