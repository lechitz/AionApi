package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var errStaleUserCacheEntry = errors.New("stale user cache entry")

// GetUserByID retrieves user profile from cache by ID.
// Returns empty User if not found (not an error).
// SECURITY: Returned user will NOT have PasswordHash - always fetch from DB for authentication.
func (s *Store) GetUserByID(ctx context.Context, userID uint64) (domain.User, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameUserGet, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "user"),
		attribute.String("user_id", strconv.FormatUint(userID, 10)),
		attribute.String("lookup_by", "id"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(UserIDKeyFormat, userID)

	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		// Cache miss is not logged as error - it's expected
		return domain.User{}, err
	}

	var dto UserCacheDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeUser, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.User{}, err
	}
	if dto.Version != UserCacheSchemaVersion {
		span.SetStatus(codes.Error, errStaleUserCacheEntry.Error())
		span.RecordError(errStaleUserCacheEntry)
		return domain.User{}, errStaleUserCacheEntry
	}

	// Convert DTO back to domain.User (PasswordHash remains empty)
	user := domain.User{
		ID:                  dto.ID,
		Name:                dto.Name,
		Username:            dto.Username,
		Email:               dto.Email,
		Locale:              dto.Locale,
		Timezone:            dto.Timezone,
		Location:            dto.Location,
		Bio:                 dto.Bio,
		AvatarURL:           dto.AvatarURL,
		OnboardingCompleted: dto.OnboardingCompleted,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
		DeletedAt:           dto.DeletedAt,
		// PasswordHash is empty - caller must fetch from DB if needed
	}

	span.SetStatus(codes.Ok, UserRetrievedSuccessfully)
	return user, nil
}

// GetUserByUsername retrieves user profile from cache by username.
// Returns empty User if not found (not an error).
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameUserGet, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "user"),
		attribute.String("username", username),
		attribute.String("lookup_by", "username"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(UserUsernameKeyFormat, username)

	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.User{}, err
	}

	var dto UserCacheDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeUser, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.User{}, err
	}
	if dto.Version != UserCacheSchemaVersion {
		span.SetStatus(codes.Error, errStaleUserCacheEntry.Error())
		span.RecordError(errStaleUserCacheEntry)
		return domain.User{}, errStaleUserCacheEntry
	}

	user := domain.User{
		ID:                  dto.ID,
		Name:                dto.Name,
		Username:            dto.Username,
		Email:               dto.Email,
		Locale:              dto.Locale,
		Timezone:            dto.Timezone,
		Location:            dto.Location,
		Bio:                 dto.Bio,
		AvatarURL:           dto.AvatarURL,
		OnboardingCompleted: dto.OnboardingCompleted,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
		DeletedAt:           dto.DeletedAt,
	}

	span.SetStatus(codes.Ok, UserRetrievedSuccessfully)
	return user, nil
}

// GetUserByEmail retrieves user profile from cache by email.
// Returns empty User if not found (not an error).
//
//nolint:dupl // Cache operations intentionally duplicated for different keys - improves maintainability
func (s *Store) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanNameUserGet, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, "user"),
		attribute.String("email", email),
		attribute.String("lookup_by", "email"),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(UserEmailKeyFormat, email)

	data, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.User{}, err
	}

	var dto UserCacheDTO
	if err := json.Unmarshal([]byte(data), &dto); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToDeserializeUser, AttributeCacheKey, cacheKey, commonkeys.Error, err)
		return domain.User{}, err
	}
	if dto.Version != UserCacheSchemaVersion {
		span.SetStatus(codes.Error, errStaleUserCacheEntry.Error())
		span.RecordError(errStaleUserCacheEntry)
		return domain.User{}, errStaleUserCacheEntry
	}

	user := domain.User{
		ID:                  dto.ID,
		Name:                dto.Name,
		Username:            dto.Username,
		Email:               dto.Email,
		Locale:              dto.Locale,
		Timezone:            dto.Timezone,
		Location:            dto.Location,
		Bio:                 dto.Bio,
		AvatarURL:           dto.AvatarURL,
		OnboardingCompleted: dto.OnboardingCompleted,
		CreatedAt:           dto.CreatedAt,
		UpdatedAt:           dto.UpdatedAt,
		DeletedAt:           dto.DeletedAt,
	}

	span.SetStatus(codes.Ok, UserRetrievedSuccessfully)
	return user, nil
}
