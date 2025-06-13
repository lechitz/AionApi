package auth

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Authenticator defines an interface for handling user authentication operations.
// Login validates user credentials and returns user details with a session token.
type Authenticator interface {
	Login(
		ctx context.Context,
		user domain.UserDomain,
		passwordReq string,
	) (domain.UserDomain, string, error)
}

// Login authenticates a user by validating credentials and generates a new token if valid. Returns the user data, token, and error if any occurs.
func (s *Service) Login(
	ctx context.Context,
	user domain.UserDomain,
	passwordReq string,
) (domain.UserDomain, string, error) {
	userDB, err := s.userRetriever.GetUserByUsername(ctx, user.Username)
	if err != nil {
		s.logger.Errorw(constants.ErrorToGetUserByUserName, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	if err := s.securityHasher.ValidatePassword(userDB.Password, passwordReq); err != nil {
		s.logger.Errorw(constants.ErrorToCompareHashAndPassword, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	tokenDomain := domain.TokenDomain{UserID: userDB.ID}

	newToken, err := s.tokenService.CreateToken(ctx, tokenDomain)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCreateToken, constants.Error, err.Error())
		return domain.UserDomain{}, "", err
	}

	s.logger.Infow(constants.SuccessToLogin, constants.UserID, userDB.ID, constants.Token, newToken)
	return userDB, newToken, nil
}
