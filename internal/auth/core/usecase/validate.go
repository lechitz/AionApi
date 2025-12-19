// Package usecase token contains use cases for managing tokens in the system.
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Validate verifies signature/exp, extracts userID from claims, and checks cache consistency.
// Returns the resolved userID and the decoded claims on success.
func (s *Service) Validate(ctx context.Context, tokenValue string) (uint64, map[string]any, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(
		ctx,
		SpanValidateToken,
		trace.WithAttributes(
			attribute.String(commonkeys.Operation, SpanValidateToken),
		),
	)
	defer span.End()

	sanitized := sanitizeTokenValue(tokenValue)

	span.AddEvent(EventVerifyToken)
	claims, err := s.authProvider.Verify(sanitized)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorInvalidToken)
		s.logger.ErrorwCtx(ctx, ErrorInvalidToken, commonkeys.Error, err.Error())
		return 0, nil, sharederrors.ErrUnauthorized(ErrorInvalidToken)
	}

	span.AddEvent(EventExtractUserID)
	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorInvalidUserIDClaim)
		s.logger.ErrorwCtx(ctx, ErrorInvalidUserIDClaim, commonkeys.Error, err.Error())
		return 0, nil, sharederrors.ErrUnauthorized(ErrorInvalidUserIDClaim)
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	span.AddEvent(EventGetTokenFromStore)
	cached, err := s.authStore.Get(ctx, userID, commonkeys.TokenTypeAccess)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToRetrieveTokenFromCache)
		s.logger.ErrorwCtx(ctx, ErrorToRetrieveTokenFromCache, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return 0, nil, errors.New(ErrorToRetrieveTokenFromCache)
	}

	span.AddEvent(EventCompareToken)
	if cached.Token == sanitized {
		// Token matches the primary (current) token
		span.AddEvent(EventTokenValidated)
		span.SetStatus(codes.Ok, SuccessTokenValidated)
		s.logger.InfowCtx(ctx, SuccessTokenValidated, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return userID, claims, nil
	}

	// Primary token doesn't match - check grace period
	s.logger.DebugwCtx(ctx, "token mismatch with primary, checking grace period",
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"provided_prefix", sanitized[:minInt(16, len(sanitized))],
		"cached_prefix", cached.Token[:minInt(16, len(cached.Token))])

	span.AddEvent("auth.token.check_grace_period")
	graceKey := buildGraceKeyForValidation(userID, sanitized)
	gracedToken, err := s.authStore.GetByKey(ctx, graceKey)
	if err == nil && gracedToken.Token == sanitized {
		// Token found in grace period and matches
		span.AddEvent("auth.token.validated_via_grace")
		span.SetStatus(codes.Ok, "token validated via grace period")
		span.SetAttributes(attribute.Bool("grace_period_used", true))
		s.logger.InfowCtx(ctx, "token validated via grace period",
			commonkeys.UserID, strconv.FormatUint(userID, 10),
			"grace_key", graceKey)
		return userID, claims, nil
	}

	// Token not found in primary or grace period - reject
	s.logger.DebugwCtx(ctx, "token not found in grace period",
		commonkeys.UserID, strconv.FormatUint(userID, 10),
		"grace_key", graceKey,
		"grace_lookup_error", err)

	span.SetStatus(codes.Error, ErrorTokenMismatch)
	s.logger.ErrorwCtx(ctx, ErrorTokenMismatch, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return 0, nil, sharederrors.ErrUnauthorized(ErrorTokenMismatch)
}

// sanitizeTokenValue strips "Bearer " prefix (case-insensitive) and trims spaces.
func sanitizeTokenValue(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 7 && strings.EqualFold(s[:7], "Bearer ") {
		return strings.TrimSpace(s[7:])
	}
	return s
}

// extractUserIDFromClaims returns the userID from claims.
// Tries claims["userID"] (canonical) then falls back to "sub".
func extractUserIDFromClaims(claims map[string]any) (uint64, error) {
	if v, ok := claims[claimskeys.UserID]; ok {
		return parseUserIDValue(v)
	}
	if v, ok := claims["sub"]; ok {
		return parseUserIDValue(v)
	}
	return 0, errors.New(ErrorInvalidUserIDClaim)
}

// parseUserIDValue converts a claim value to uint64 safely.
func parseUserIDValue(v any) (uint64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseUint(strings.TrimSpace(t), 10, 64)
	case json.Number:
		return strconv.ParseUint(t.String(), 10, 64)
	case float64:
		if t < 0 || t != math.Trunc(t) {
			return 0, errors.New(ErrorInvalidUserIDClaim)
		}
		return uint64(t), nil
	default:
		return 0, errors.New(ErrorInvalidUserIDClaim)
	}
}

// buildGraceKeyForValidation constructs a Redis key for retrieving a token from grace period.
// Format: auth:grace:{userID}:{tokenHash}.
func buildGraceKeyForValidation(userID uint64, token string) string {
	return fmt.Sprintf("auth:grace:%d:%s", userID, hashTokenForValidation(token))
}

// hashTokenForValidation creates a short hash of the token for use in Redis keys.
// Uses first 16 + last 16 characters to keep keys short while maintaining uniqueness.
func hashTokenForValidation(token string) string {
	if len(token) < 32 {
		return token
	}
	return token[:16] + token[len(token)-16:]
}

// minInt returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
