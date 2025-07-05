// Package cache provides methods for interacting with a Redis-based token storage.
package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/constants"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TokenRepository provides methods to interact with a Redis-based token storage.
type TokenRepository struct {
	cache  *redis.Client
	logger output.Logger
}

// NewTokenRepository creates and returns a new instance of TokenRepository with the provided Redis client and logger.
func NewTokenRepository(cache *redis.Client, logger output.Logger) *TokenRepository {
	return &TokenRepository{
		cache:  cache,
		logger: logger,
	}
}

// Save stores a token in the Redis cache with a 24-hour expiration time and logs errors if the operation fails.
func (t *TokenRepository) Save(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer("TokenRepository")
	ctx, span := tr.Start(ctx, "Save Token", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "save"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour

	if err := t.cache.Set(ctx, key, token.Token, expiration).Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		t.logger.Errorw(constants.ErrorToSaveTokenToRedis, constants.Key, key, constants.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, "token saved successfully")
	return nil
}

// Get retrieves a token associated with a user ID from the Redis cache or returns an error if the token is not found or another issue occurs.
func (t *TokenRepository) Get(ctx context.Context, token domain.TokenDomain) (string, error) {
	tr := otel.Tracer("TokenRepository")
	ctx, span := tr.Start(ctx, "Get Token", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "get"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)

	value, err := t.cache.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) || err.Error() == "redis: nil" {
			span.SetStatus(codes.Ok, "token not found (business as usual)")
			return "", nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		t.logger.Errorw(
			constants.ErrorToGetTokenFromRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return "", err
	}

	span.SetStatus(codes.Ok, "token retrieved successfully")
	return value, nil
}

// Update updates an existing token in the Redis cache with a 24-hour expiration and logs success or failure.
func (t *TokenRepository) Update(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer("TokenRepository")
	ctx, span := tr.Start(ctx, "Update Token", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "update"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour

	if err := t.cache.Set(ctx, key, token.Token, expiration).Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		t.logger.Errorw(
			constants.ErrorToUpdateTokenInRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return err
	}

	span.SetStatus(codes.Ok, "token updated successfully")
	t.logger.Infow(constants.SuccessToUpdateTokenInRedis, constants.Key, key)
	return nil
}

// Delete removes a token associated with a user ID from the Redis cache and logs any errors if the operation fails.
func (t *TokenRepository) Delete(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer("TokenRepository")
	ctx, span := tr.Start(ctx, "Delete Token", trace.WithAttributes(
		attribute.String(constants.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "delete"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)

	if err := t.cache.Del(ctx, key).Err(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		t.logger.Errorw(
			constants.ErrorToDeleteTokenFromRedis,
			constants.Key,
			key,
			constants.Error,
			err,
		)
		return err
	}

	span.SetStatus(codes.Ok, "token deleted successfully")
	return nil
}

// formatTokenKey generates a Redis key for storing a user token by appending the user ID to a predefined base string.
func (t *TokenRepository) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token_user_%d", userID)
}
