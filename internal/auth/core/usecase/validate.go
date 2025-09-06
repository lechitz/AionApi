// Package usecase token contains use cases for managing tokens in the system.
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

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
	claims, err := s.tokenProvider.Verify(ctx, sanitized)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorInvalidToken)
		s.logger.ErrorwCtx(ctx, ErrorInvalidToken, commonkeys.Error, err.Error())
		return 0, nil, errors.New(ErrorInvalidToken)
	}

	span.AddEvent(EventExtractUserID)
	userID, err := extractUserIDFromClaims(claims)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorInvalidUserIDClaim)
		s.logger.ErrorwCtx(ctx, ErrorInvalidUserIDClaim, commonkeys.Error, err.Error())
		return 0, nil, errors.New(ErrorInvalidUserIDClaim)
	}
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)))

	span.AddEvent(EventGetTokenFromStore)
	cached, err := s.tokenStore.Get(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToRetrieveTokenFromCache)
		s.logger.ErrorwCtx(ctx, ErrorToRetrieveTokenFromCache, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(userID, 10))
		return 0, nil, errors.New(ErrorToRetrieveTokenFromCache)
	}

	span.AddEvent(EventCompareToken)
	if cached.Token == "" || cached.Token != sanitized {
		span.SetStatus(codes.Error, ErrorTokenMismatch)
		s.logger.ErrorwCtx(ctx, ErrorTokenMismatch, commonkeys.UserID, strconv.FormatUint(userID, 10))
		return 0, nil, errors.New(ErrorTokenMismatch)
	}

	span.AddEvent(EventTokenValidated)
	span.SetStatus(codes.Ok, SuccessTokenValidated)
	s.logger.InfowCtx(ctx, SuccessTokenValidated, commonkeys.UserID, strconv.FormatUint(userID, 10))
	return userID, claims, nil
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

// parseUserIDValue converts a claim value to uint64 safely (no float64 precision loss).
func parseUserIDValue(v any) (uint64, error) {
	switch t := v.(type) {
	case string:
		return strconv.ParseUint(strings.TrimSpace(t), 10, 64)
	case json.Number:
		return strconv.ParseUint(t.String(), 10, 64)
	default:
		return 0, errors.New(ErrorInvalidUserIDClaim)
	}
}
