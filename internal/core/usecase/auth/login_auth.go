package auth

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

type Authenticator interface {
	Login(ctx context.Context, user domain.UserDomain, passwordReq string) (domain.UserDomain, string, error)
}

func (s *AuthService) Login(
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
