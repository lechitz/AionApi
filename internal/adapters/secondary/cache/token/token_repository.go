// Package token provides methods for managing tokens in the cache.Cache.
package token

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/shared/commonkeys"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/token/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Repository is a repository for managing tokens.
type Repository struct {
	cache  output.Cache
	logger output.Logger
}

// NewTokenRepository creates a new instance of Repository with a given cache.Cache and logger.Logger.
func NewTokenRepository(cache output.Cache, logger output.Logger) *Repository {
	return &Repository{
		cache:  cache,
		logger: logger,
	}
}

// TODO: Ajustar magic strings.

// Save saves a token in the cache.Cache for the specified userID. Returns an error if the operation fails.
func (t *Repository) Save(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer("Repository")
	ctx, span := tr.Start(ctx, "Save Token", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "save"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)
	expiration := 24 * time.Hour

	if err := t.cache.Set(ctx, key, token.Token, expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		// TODO: ajustar o uso de TokenUser no logger.

		t.logger.Errorw(constants.ErrorToSaveTokenToRedis, commonkeys.TokenUser, key, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, "token saved successfully")
	return nil
}

// Get retrieves the token associated with the specified user ID from the cache. Returns an empty string and nil if the token does not exist.
func (t *Repository) Get(ctx context.Context, token domain.TokenDomain) (string, error) {
	tr := otel.Tracer("Repository")
	ctx, span := tr.Start(ctx, "Get Token", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "get"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)

	value, err := t.cache.Get(ctx, key)
	if err != nil {
		if errors.Is(err, output.ErrNil) {
			span.SetStatus(codes.Ok, "token not found (business as usual)")
			return "", nil
		}

		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		// TODO: ajustar o uso de TokenUser no logger.

		t.logger.Errorw(constants.ErrorToGetTokenFromRedis, commonkeys.TokenUser, key, constants.Error, err)
		return "", err
	}

	span.SetStatus(codes.Ok, "token retrieved successfully")
	return value, nil
}

// Delete deletes the token associated with the specified user ID from the cache. Returns an error if the operation fails.
func (t *Repository) Delete(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer("Repository")
	ctx, span := tr.Start(ctx, "Delete Token", trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String("operation", "delete"),
	))
	defer span.End()

	key := t.formatTokenKey(token.UserID)

	if err := t.cache.Del(ctx, key); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)

		// TODO: ajustar o uso de TokenUser no logger.

		t.logger.Errorw(constants.ErrorToDeleteTokenFromRedis, commonkeys.TokenUser, key, constants.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, "token deleted successfully")
	return nil
}

// formatTokenKey formats the token key for the specified userID.
func (t *Repository) formatTokenKey(userID uint64) string {
	return fmt.Sprintf("token:user:%d", userID)
}
