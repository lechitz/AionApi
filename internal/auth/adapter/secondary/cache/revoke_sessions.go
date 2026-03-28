package cache

import (
	"context"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
)

// RevokeUserSessions removes access and refresh tokens for the user.
func (s *Store) RevokeUserSessions(ctx context.Context, userID uint64) error {
	if err := s.Delete(ctx, userID, commonkeys.TokenTypeAccess); err != nil {
		return err
	}
	if err := s.Delete(ctx, userID, commonkeys.TokenTypeRefresh); err != nil {
		return err
	}
	return nil
}
