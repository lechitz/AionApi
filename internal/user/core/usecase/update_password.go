package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	userdomain "github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UpdatePassword updates a user's password after validating the old password and hashing the new password, then returns the updated user and a new token.
func (s *Service) UpdatePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) (string, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanUpdateUserPassword)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanUpdateUserPassword),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	user, err := s.getUserWithCache(ctx, span, userID)
	if err != nil {
		return "", err
	}

	hashedPassword, err := s.validateAndHashPassword(ctx, span, user, oldPassword, newPassword)
	if err != nil {
		return "", err
	}

	fields := map[string]interface{}{
		commonkeys.Password:      hashedPassword,
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	updatedUser, err := s.userRepository.Update(ctx, user.ID, fields)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToUpdatePassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToUpdatePassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%w: %w", ErrUpdatePassword, err)
	}

	tokenValue, err := s.generateAndSaveRefreshToken(ctx, span, updatedUser.ID)
	if err != nil {
		return "", err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return tokenValue, nil
}

// getUserWithCache attempts to fetch user from cache first, falling back to database if needed.
// This helper maintains the cache-first pattern and proper instrumentation.
func (s *Service) getUserWithCache(ctx context.Context, span trace.Span, userID uint64) (userdomain.User, error) {
	span.AddEvent(SpanEventCheckCache)
	user, err := s.userCache.GetUserByID(ctx, userID)
	if err != nil || user.ID == 0 {
		span.AddEvent(SpanEventCacheMiss)
		user, err = s.userRepository.GetByID(ctx, userID)
		if err != nil {
			span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetSelf))
			span.RecordError(err)
			s.logger.ErrorwCtx(ctx, ErrorToGetSelf, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
			return userdomain.User{}, fmt.Errorf("%s: %w", ErrorToGetSelf, err)
		}
		return user, nil
	}

	span.AddEvent(SpanEventCacheHit)
	s.logger.InfowCtx(ctx, InfoUserRetrievedFromCache, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return user, nil
}

// validateAndHashPassword verifies the old password and generates hash for new password.
// Returns the hashed password or an error if validation/hashing fails.
func (s *Service) validateAndHashPassword(ctx context.Context, span trace.Span, user userdomain.User, oldPassword, newPassword string) (string, error) {
	if err := s.hasher.Compare(user.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%w: %w", ErrCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToHashPassword, err)
	}

	return hashedPassword, nil
}

// generateAndSaveRefreshToken creates a new refresh token and saves it to the auth store.
// The token is saved with TTL extracted from JWT claims to ensure automatic expiration.
func (s *Service) generateAndSaveRefreshToken(ctx context.Context, span trace.Span, userID uint64) (string, error) {
	tokenValue, err := s.tokenProvider.GenerateRefreshToken(userID)
	if err != nil || tokenValue == "" {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		if err != nil {
			span.RecordError(err)
			s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
			return "", fmt.Errorf("%w: %w", ErrCreateToken, err)
		}
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return "", fmt.Errorf("%w: empty token", ErrCreateToken)
	}

	// Compute TTL from token claims to sync Redis expiration with JWT expiration
	expiration := time.Duration(0)
	if claims, err := s.tokenProvider.Verify(tokenValue); err == nil {
		expiration = claimsTTLFromVerifyResult(claims)
	}

	if err := s.authStore.Save(ctx, domain.NewAccessToken(tokenValue, userID).ToAuth(), expiration); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return "", fmt.Errorf("%w: %w", ErrCreateToken, err)
	}

	return tokenValue, nil
}

// claimsTTLFromVerifyResult computes the TTL (Time To Live) from JWT claims 'exp' value.
// This ensures tokens expire in Redis at the same time as the JWT itself.
// Returns 0 if the expiration claim is missing or invalid.
func claimsTTLFromVerifyResult(claims map[string]any) time.Duration {
	if claims == nil {
		return 0
	}

	expirationValue, exists := claims[claimskeys.Exp]
	if !exists {
		return 0
	}

	switch typedValue := expirationValue.(type) {
	case float64:
		expirationUnix := int64(typedValue)
		return time.Until(time.Unix(expirationUnix, 0))
	case int64:
		return time.Until(time.Unix(typedValue, 0))
	case int:
		return time.Until(time.Unix(int64(typedValue), 0))
	case string:
		expirationUnix, err := strconv.ParseInt(typedValue, 10, 64)
		if err != nil {
			return 0
		}
		return time.Until(time.Unix(expirationUnix, 0))
	case json.Number:
		expirationUnix, err := typedValue.Int64()
		if err != nil {
			return 0
		}
		return time.Until(time.Unix(expirationUnix, 0))
	default:
		return 0
	}
}
