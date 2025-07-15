// Package token provides methods for managing tokens in the cache.Cache.
package token

import (
	"context"
	"fmt"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/cache/token/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Store is a repository for managing tokens.
type Store struct {
	cache  output.Cache
	logger output.ContextLogger
}

// NewStore creates a new instance of Store with a given cache.Cache and logger.Logger.
func NewStore(cache output.Cache, logger output.ContextLogger) *Store {
	return &Store{
		cache:  cache,
		logger: logger,
	}
}

// Save saves a token in the cache.Cache for the specified userID. Returns an error if the operation fails.
func (s *Store) Save(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer(constants.SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, constants.SpanNameTokenSave, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationSave),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	key := s.formatTokenKey(token.UserID)
	expiration := constants.TokenExpirationDefault

	if err := s.cache.Set(ctx, key, token.Key, expiration); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(constants.ErrorToSaveTokenToRedis, commonkeys.Token, key, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, constants.TokenSavedSuccessfully)
	return nil
}

// Get retrieves the token associated with the specified user ID from the cache. Returns an empty string and nil if the token does not exist.
func (s *Store) Get(ctx context.Context, userID uint64) (string, error) {
	tr := otel.Tracer(constants.SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, constants.SpanNameTokenGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationGet),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	tokenKey := s.formatTokenKey(userID)

	token, err := s.cache.Get(ctx, tokenKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(constants.ErrorToGetTokenFromRedis, commonkeys.Token, token, commonkeys.Error, err)
		return "", err
	}

	if token == "" {
		span.SetStatus(codes.Ok, sharederrors.ErrTokenNotFound)
		return "", nil
	}

	span.SetStatus(codes.Ok, constants.TokenRetrievedSuccessfully)
	return token, nil
}

// Delete deletes the token associated with the specified user ID from the cache. Returns an error if the operation fails.
func (s *Store) Delete(ctx context.Context, token domain.TokenDomain) error {
	tr := otel.Tracer(constants.SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, constants.SpanNameTokenDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(token.UserID, 10)),
		attribute.String(commonkeys.Operation, commonkeys.OperationDelete),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	key := s.formatTokenKey(token.UserID)

	if err := s.cache.Del(ctx, key); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(constants.ErrorToDeleteTokenFromRedis, commonkeys.Token, key, commonkeys.Error, err)
		return err
	}

	span.SetStatus(codes.Ok, constants.TokenDeletedSuccessfully)
	return nil
}

// formatTokenKey formats the token key for the specified userID.
func (s *Store) formatTokenKey(userID uint64) string {
	return fmt.Sprintf(constants.TokenUserKeyFormat, userID)
}
