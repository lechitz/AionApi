package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	authDomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// RefreshTokenRenewal renews an access token using a valid refresh token string.
func (s *Service) RefreshTokenRenewal(ctx context.Context, refreshTokenValue string) (string, string, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanRefreshTokenRenewal)
	defer span.End()

	span.SetAttributes(attribute.String(commonkeys.Operation, OperationRefreshTokenRenewal))

	span.AddEvent(EventValidateRefreshToken)
	userID, err := s.validateRefreshTokenAndStored(ctx, refreshTokenValue)
	if err != nil {
		s.handleInvalidRefreshToken(ctx, span, err)
		return "", "", ErrInvalidRefreshToken
	}

	span.AddEvent(EventGenerateAndSaveNewTokens)
	newAccessToken, newRefreshToken, err := s.generateAndSaveTokens(ctx, userID)
	if err != nil {
		s.handleCreateTokenError(ctx, span, err)
		return "", "", ErrTokenCreation
	}

	span.SetStatus(codes.Ok, SuccessRefreshTokenRenewed)
	s.logger.InfowCtx(ctx, SuccessRefreshTokenRenewed, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return newAccessToken, newRefreshToken, nil
}

// validateRefreshTokenAndStored verifies the refresh token, extracts the userID and
// ensures the token matches what's stored in the auth store.
func (s *Service) validateRefreshTokenAndStored(ctx context.Context, refreshTokenValue string) (uint64, error) {
	claims, err := s.authProvider.Verify(refreshTokenValue)
	if err != nil {
		return 0, err
	}
	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		return 0, err
	}
	stored, err := s.authStore.Get(ctx, userID, commonkeys.TokenTypeRefresh)
	if err != nil || stored.Token != refreshTokenValue {
		return 0, ErrInvalidRefreshToken
	}
	return userID, nil
}

// generateAndSaveTokens creates access and refresh tokens for the user and uses
// saveAuthIfTTL to persist them when TTLs are available.
// It also moves the old access token to a grace period cache to prevent race conditions
// in multi-tab scenarios.
func (s *Service) generateAndSaveTokens(ctx context.Context, userID uint64) (string, string, error) {
	// Retrieve the old token before generating new ones (for grace period)
	// Note: Get() may return empty token (not an error) or actual error
	oldAuth, _ := s.authStore.Get(ctx, userID, commonkeys.TokenTypeAccess)
	if oldAuth.Token == "" {
		// First login or token not found - not an error, just log debug
		s.logger.DebugwCtx(ctx, LogNoPreviousTokenForGrace, commonkeys.UserID, userID)
	}

	user, err := s.userCache.GetUserByID(ctx, userID)
	if err != nil {
		user, err = s.userRepository.GetByID(ctx, userID)
		if err != nil {
			return "", "", fmt.Errorf(ErrorToGetUserData+": %w", err)
		}

		err = s.userCache.SaveUser(ctx, user, 0)
		if err != nil {
			s.logger.WarnwCtx(ctx, LogFailedToCacheUserData,
				commonkeys.UserID, userID,
				commonkeys.Error, err.Error(),
			)
		}
	}

	roles, err := s.getRolesWithCache(ctx, userID)
	if err != nil {
		return "", "", err
	}
	claimsForAccess := map[string]any{
		claimskeys.UserID:   userID,
		claimskeys.Username: user.Username,
		claimskeys.Email:    user.Email,
		claimskeys.Name:     user.Name,
		claimskeys.Roles:    roles,
	}

	span := trace.SpanFromContext(ctx)
	span.AddEvent(EventGenerateToken)
	access, err := s.authProvider.GenerateAccessToken(userID, claimsForAccess)
	if err != nil {
		return "", "", err
	}

	// Save new access token as primary
	span.AddEvent(EventSaveAccessTokenToStore)
	if err := s.saveAuthIfTTL(ctx, authDomain.NewAccessToken(access, userID).ToAuth(), access); err != nil {
		return "", "", err
	}

	// Move old token to grace period if it exists and is different from new token
	shouldSkipGrace := oldAuth.Token == "" || oldAuth.Token == access
	if shouldSkipGrace {
		// Log why we didn't create grace period
		if oldAuth.Token == "" {
			s.logger.DebugwCtx(ctx, LogSkippingGraceNoOldToken, commonkeys.UserID, userID)
		} else {
			s.logger.DebugwCtx(ctx, LogSkippingGraceTokensIdentical, commonkeys.UserID, userID)
		}
	} else {
		s.moveTokenToGracePeriod(ctx, userID, oldAuth, access)
	}

	span.AddEvent(EventGenerateRefreshToken)
	refresh, err := s.authProvider.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
	span.AddEvent(EventSaveRefreshTokenToStore)
	if err := s.saveAuthIfTTL(ctx, authDomain.NewRefreshToken(refresh, userID).ToAuth(), refresh); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

// tokenTTLFromValue returns the TTL derived from the token's 'exp' claim or 0 when unavailable.
func (s *Service) tokenTTLFromValue(token string) time.Duration {
	if token == "" {
		return 0
	}
	if c, err := s.authProvider.Verify(token); err == nil {
		return claimsTTLFromVerifyResult(c)
	}
	return 0
}

// saveAuthIfTTL saves the auth entry if its computed TTL is positive; otherwise does nothing.
func (s *Service) saveAuthIfTTL(ctx context.Context, a authDomain.Auth, token string) error {
	ttl := s.tokenTTLFromValue(token)
	if ttl > 0 {
		return s.authStore.Save(ctx, a, ttl)
	}
	return nil
}

// handleInvalidRefreshToken records the error in span and logs it.
func (s *Service) handleInvalidRefreshToken(ctx context.Context, span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, ErrorInvalidRefreshToken)
	s.logger.ErrorwCtx(ctx, ErrorInvalidRefreshToken, commonkeys.Error, err)
}

// handleCreateTokenError records token creation errors in span and logs them.
func (s *Service) handleCreateTokenError(ctx context.Context, span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, ErrorToCreateToken)
	s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err)
}

// moveTokenToGracePeriod moves the old token to grace period storage.
func (s *Service) moveTokenToGracePeriod(ctx context.Context, userID uint64, oldAuth authDomain.Auth, newToken string) {
	graceKey := buildGraceKey(userID, oldAuth.Token)

	// Safely get token prefixes for logging
	oldPrefix := oldAuth.Token
	if len(oldAuth.Token) > 16 {
		oldPrefix = oldAuth.Token[:16]
	}
	newPrefix := newToken
	if len(newToken) > 16 {
		newPrefix = newToken[:16]
	}

	s.logger.DebugwCtx(ctx, LogMovingTokenToGrace,
		commonkeys.UserID, userID,
		commonkeys.GraceKey, graceKey,
		commonkeys.OldTokenPrefix, oldPrefix,
		commonkeys.NewTokenPrefix, newPrefix)

	if err := s.authStore.SaveWithKey(ctx, graceKey, oldAuth, GracePeriodDuration); err != nil {
		// Log error but don't fail the refresh - grace period is a best-effort feature
		s.logger.ErrorwCtx(ctx, LogFailedToSaveTokenToGrace,
			commonkeys.UserID, userID,
			commonkeys.Error, err)
	} else {
		s.logger.InfowCtx(ctx, LogTokenMovedToGraceSuccess,
			commonkeys.UserID, userID,
			commonkeys.GraceKey, graceKey,
			commonkeys.GraceTTL, GracePeriodDuration.String())
	}
}

// buildGraceKey constructs a Redis key for storing a token in grace period.
// Format: auth:grace:{userID}:{tokenHash}.
func buildGraceKey(userID uint64, token string) string {
	return fmt.Sprintf(AuthGraceKeyPrefix+":%d:%s", userID, hashToken(token))
}

// hashToken creates a short hash of the token for use in Redis keys.
// Uses first 16 + last 16 characters to keep keys short while maintaining uniqueness.
func hashToken(token string) string {
	if len(token) < 32 {
		return token
	}
	return token[:16] + token[len(token)-16:]
}
