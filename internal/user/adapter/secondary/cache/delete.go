package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// DeleteUser removes user profile from cache (all keys: ID, username, email).
func (s *Store) DeleteUser(ctx context.Context, userID uint64, username, email string) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameUserDelete, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationDelete),
		attribute.String(commonkeys.Entity, "user"),
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
	))
	defer span.End()

	// Delete by ID (primary)
	cacheKeyID := fmt.Sprintf(UserIDKeyFormat, userID)
	if err := s.cache.Del(ctx, cacheKeyID); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeleteUserFromCache, AttributeCacheKey, cacheKeyID, commonkeys.Error, err)
		return err
	}

	// Delete by username (secondary - best effort)
	if username != "" {
		cacheKeyUsername := fmt.Sprintf(UserUsernameKeyFormat, username)
		if err := s.cache.Del(ctx, cacheKeyUsername); err != nil {
			s.logger.Warnw("failed to delete user from cache by username",
				AttributeCacheKey, cacheKeyUsername,
				commonkeys.Error, err,
			)
		}
	}

	// Delete by email (secondary - best effort)
	if email != "" {
		cacheKeyEmail := fmt.Sprintf(UserEmailKeyFormat, email)
		if err := s.cache.Del(ctx, cacheKeyEmail); err != nil {
			s.logger.Warnw("failed to delete user from cache by email",
				AttributeCacheKey, cacheKeyEmail,
				commonkeys.Error, err,
			)
		}
	}

	span.SetStatus(codes.Ok, UserDeletedSuccessfully)
	return nil
}
