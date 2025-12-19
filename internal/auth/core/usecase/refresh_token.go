package usecase

import (
	"context"
	"errors"
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
	ctx, span := tracer.Start(ctx, "RefreshTokenRenewal")
	defer span.End()

	span.SetAttributes(attribute.String(commonkeys.Operation, "RefreshTokenRenewal"))

	// validate incoming refresh token and stored value
	userID, err := s.validateRefreshTokenAndStored(ctx, refreshTokenValue)
	if err != nil {
		s.handleInvalidRefreshToken(ctx, span, err)
		return "", "", errors.New("invalid refresh token")
	}

	// generate new tokens and persist them if they have TTL
	newAccessToken, newRefreshToken, err := s.generateAndSaveTokens(ctx, userID)
	if err != nil {
		s.handleCreateTokenError(ctx, span, err)
		return "", "", errors.New(ErrorToCreateToken)
	}

	span.SetStatus(codes.Ok, "refresh token renewed successfully")
	s.logger.InfowCtx(ctx, "refresh token renewed successfully", commonkeys.UserID, strconv.FormatUint(userID, 10))
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
		return 0, errors.New("invalid refresh token")
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
		s.logger.DebugwCtx(ctx, "no previous token found for grace period", commonkeys.UserID, userID)
	}

	// Generate new access token
	claimsForAccess := map[string]any{claimskeys.UserID: userID}
	access, err := s.authProvider.GenerateAccessToken(userID, claimsForAccess)
	if err != nil {
		return "", "", err
	}

	// Save new access token as primary
	if err := s.saveAuthIfTTL(ctx, authDomain.NewAccessToken(access, userID).ToAuth(), access); err != nil {
		return "", "", err
	}

	// Move old token to grace period if it exists and is different from new token
	if oldAuth.Token != "" && oldAuth.Token != access {
		graceKey := buildGraceKey(userID, oldAuth.Token)
		s.logger.DebugwCtx(ctx, "moving token to grace period",
			commonkeys.UserID, userID,
			"grace_key", graceKey,
			"old_token_prefix", oldAuth.Token[:min(16, len(oldAuth.Token))],
			"new_token_prefix", access[:min(16, len(access))])

		if err := s.authStore.SaveWithKey(ctx, graceKey, oldAuth, GracePeriodDuration); err != nil {
			// Log error but don't fail the refresh - grace period is a best-effort feature
			s.logger.ErrorwCtx(ctx, "failed to save token to grace period",
				commonkeys.UserID, userID,
				commonkeys.Error, err)
		} else {
			s.logger.InfowCtx(ctx, "token moved to grace period successfully",
				commonkeys.UserID, userID,
				"grace_key", graceKey,
				"grace_ttl", GracePeriodDuration.String())
		}
	} else {
		// Log why we didn't create grace period
		if oldAuth.Token == "" {
			s.logger.DebugwCtx(ctx, "skipping grace period: no old token", commonkeys.UserID, userID)
		} else if oldAuth.Token == access {
			s.logger.DebugwCtx(ctx, "skipping grace period: tokens are identical", commonkeys.UserID, userID)
		}
	}

	// Generate and save new refresh token
	refresh, err := s.authProvider.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}
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
	span.SetStatus(codes.Error, "invalid refresh token")
	s.logger.ErrorwCtx(ctx, "invalid refresh token", commonkeys.Error, err)
}

// handleCreateTokenError records token creation errors in span and logs them.
func (s *Service) handleCreateTokenError(ctx context.Context, span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, ErrorToCreateToken)
	s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err)
}

// buildGraceKey constructs a Redis key for storing a token in grace period.
// Format: auth:grace:{userID}:{tokenHash}.
func buildGraceKey(userID uint64, token string) string {
	return fmt.Sprintf("auth:grace:%d:%s", userID, hashToken(token))
}

// hashToken creates a short hash of the token for use in Redis keys.
// Uses first 16 + last 16 characters to keep keys short while maintaining uniqueness.
func hashToken(token string) string {
	if len(token) < 32 {
		return token
	}
	return token[:16] + token[len(token)-16:]
}

// min returns the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
