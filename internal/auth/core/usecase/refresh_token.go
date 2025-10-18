package usecase

import (
	"context"
	"errors"
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
func (s *Service) generateAndSaveTokens(ctx context.Context, userID uint64) (string, string, error) {
	claimsForAccess := map[string]any{claimskeys.UserID: userID}
	access, err := s.authProvider.GenerateAccessToken(userID, claimsForAccess)
	if err != nil {
		return "", "", err
	}
	if err := s.saveAuthIfTTL(ctx, authDomain.NewAccessToken(access, userID).ToAuth(), access); err != nil {
		return "", "", err
	}
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
