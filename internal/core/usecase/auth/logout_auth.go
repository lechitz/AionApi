package auth

import (
	"context"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/def"

	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Logout revokes a user's authentication token, effectively logging them out. Returns an error if token verification or deletion fails.
func (s *Service) Logout(ctx context.Context, token string) error {
	userID, _, err := s.tokenService.VerifyToken(ctx, token)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCheckToken, def.Error, err.Error())
		return err
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToRevokeToken, def.Error, err.Error(), def.CtxUserID, userID)
		return err
	}

	s.logger.Infow(constants.SuccessUserLoggedOut, def.CtxUserID, userID)

	return nil
}
