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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToGetSelf))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToGetSelf, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToGetSelf, err)
	}

	if err := s.hasher.Compare(user.Password, oldPassword); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToCompareHashAndPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCompareHashAndPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToCompareHashAndPassword, err)
	}

	hashedPassword, err := s.hasher.Hash(newPassword)
	if err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, ErrorToHashPassword))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToHashPassword, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(user.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToHashPassword, err)
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
		return "", fmt.Errorf("%s: %w", ErrorToUpdatePassword, err)
	}

	tokenValue, err := s.tokenProvider.GenerateRefreshToken(updatedUser.ID)
	if err != nil || tokenValue == "" {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		if err != nil {
			span.RecordError(err)
			s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
			return "", fmt.Errorf("%s: %w", ErrorToCreateToken, err)
		}
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: empty token", ErrorToCreateToken)
	}

	// compute TTL from tokenValue via Verify
	expiration := time.Duration(0)
	if c, err := s.tokenProvider.Verify(tokenValue); err == nil {
		expiration = claimsTTLFromVerifyResult(c)
	}

	if err := s.authStore.Save(ctx, domain.NewAccessToken(tokenValue, updatedUser.ID).ToAuth(), expiration); err != nil {
		span.SetAttributes(attribute.String(commonkeys.Status, sharederrors.ErrMsgCreateToken))
		span.RecordError(err)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error(), commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))
		return "", fmt.Errorf("%s: %w", ErrorToCreateToken, err)
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, commonkeys.StatusSuccess),
		attribute.String(commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10)),
	)
	s.logger.InfowCtx(ctx, SuccessPasswordUpdated, commonkeys.UserID, strconv.FormatUint(updatedUser.ID, 10))

	return tokenValue, nil
}

// helper to compute TTL from claims 'exp' value (copied from auth usecase).
func claimsTTLFromVerifyResult(claims map[string]any) time.Duration {
	if claims == nil {
		return 0
	}
	v, ok := claims[claimskeys.Exp]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case float64:
		exp := int64(x)
		return time.Until(time.Unix(exp, 0))
	case int64:
		return time.Until(time.Unix(x, 0))
	case int:
		return time.Until(time.Unix(int64(x), 0))
	case string:
		n, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return 0
		}
		return time.Until(time.Unix(n, 0))
	case json.Number:
		n, err := x.Int64()
		if err != nil {
			return 0
		}
		return time.Until(time.Unix(n, 0))
	default:
		return 0
	}
}
