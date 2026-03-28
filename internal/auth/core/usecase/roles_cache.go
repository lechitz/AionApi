package usecase

import (
	"context"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel/trace"
)

func (s *Service) getRolesWithCache(ctx context.Context, userID uint64) ([]string, error) {
	span := trace.SpanFromContext(ctx)

	if s.roleCache != nil {
		span.AddEvent(EventGetRolesFromCache)
		roles, err := s.roleCache.GetRoles(ctx, userID)
		if err == nil && len(roles) > 0 {
			return roles, nil
		}
		if err != nil {
			s.logger.WarnwCtx(ctx, WarnFailedToGetRolesFromCache,
				commonkeys.UserID, userID,
				commonkeys.Error, err,
			)
		}
	}

	roles, err := s.rolesReader.GetRolesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if s.roleCache != nil {
		span.AddEvent(EventSaveRolesToCache)
		if err := s.roleCache.SaveRoles(ctx, userID, roles, 0); err != nil {
			s.logger.WarnwCtx(ctx, WarnFailedToCacheRoles,
				commonkeys.UserID, userID,
				commonkeys.Error, err,
			)
		}
	}

	return roles, nil
}
