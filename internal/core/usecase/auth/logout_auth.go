package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
)

// Logout revokes a user's authentication token, effectively logging them out. Returns an error if token verification or deletion fails.
func (s *Service) Logout(ctx context.Context, token string) error {
	userID, _, err := s.tokenService.GetToken(ctx, token)
	if err != nil {
		s.logger.Errorw(constants.ErrorToCheckToken, commonkeys.Error, err.Error())
		return fmt.Errorf("%s: %w", constants.ErrorToCheckToken, err)
	}

	tokenDomain := domain.TokenDomain{
		UserID: userID,
		Token:  token,
	}

	if err := s.tokenService.Delete(ctx, tokenDomain); err != nil {
		s.logger.Errorw(constants.ErrorToRevokeToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return fmt.Errorf("%s: %w", constants.ErrorToRevokeToken, err)
	}

	s.logger.Infow(constants.SuccessUserLoggedOut, commonkeys.UserID, strconv.FormatUint(userID, 10))

	return nil
}
