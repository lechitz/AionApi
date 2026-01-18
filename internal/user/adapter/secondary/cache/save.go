package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SaveUser persists user profile data in cache (WITHOUT password hash).
// SECURITY: Password hash is intentionally excluded from cache.
func (s *Store) SaveUser(ctx context.Context, user domain.User, expiration time.Duration) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameUserSave, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationSave),
		attribute.String(commonkeys.Entity, commonkeys.User),
		attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)),
	))
	defer span.End()

	if expiration <= 0 {
		expiration = UserExpirationDefault
	}

	// Convert domain.User to UserCacheDTO (excludes PasswordHash)
	dto := UserCacheDTO{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}

	data, err := json.Marshal(dto)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSerializeUser, "user_id", user.ID, commonkeys.Error, err)
		return err
	}

	// Save by ID (primary key)
	cacheKeyID := fmt.Sprintf(UserIDKeyFormat, user.ID)
	if err := s.cache.Set(ctx, cacheKeyID, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToSaveUserToCache, AttributeCacheKey, cacheKeyID, commonkeys.Error, err)
		return err
	}

	// Also save by username for username-based lookups
	cacheKeyUsername := fmt.Sprintf(UserUsernameKeyFormat, user.Username)
	if err := s.cache.Set(ctx, cacheKeyUsername, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Warnw("failed to save user to cache by username",
			AttributeCacheKey, cacheKeyUsername,
			commonkeys.Error, err,
		)
		// Don't return error - primary key cache succeeded
	}

	// Also save by email for email-based lookups
	cacheKeyEmail := fmt.Sprintf(UserEmailKeyFormat, user.Email)
	if err := s.cache.Set(ctx, cacheKeyEmail, string(data), expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Warnw("failed to save user to cache by email",
			AttributeCacheKey, cacheKeyEmail,
			commonkeys.Error, err,
		)
		// Don't return error - primary key cache succeeded
	}

	span.SetStatus(codes.Ok, UserSavedSuccessfully)
	return nil
}
