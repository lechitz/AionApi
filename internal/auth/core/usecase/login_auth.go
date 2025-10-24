// Package usecase (auth) contains use cases for authenticating users and generating tokens.
package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	authDomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/claimskeys"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	userDomain "github.com/lechitz/AionApi/internal/user/core/domain"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
func (s *Service) Login(ctx context.Context, usernameReq, passwordReq string) (userDomain.User, string, string, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanLogin)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, SpanLogin),
		attribute.String(commonkeys.Username, usernameReq),
	)

	span.AddEvent(EventLookupUser)
	user, err := s.userRepository.GetByUsername(ctx, usernameReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToGetUserByUserName)
		s.logger.ErrorwCtx(ctx, ErrorToGetUserByUserName, commonkeys.Error, err.Error())
		return userDomain.User{}, "", "", errors.New(ErrorToGetUserByUserName)
	}
	if user.ID == 0 {
		span.SetStatus(codes.Error, UserNotFoundOrInvalidCredentials)
		s.logger.WarnwCtx(ctx, UserNotFoundOrInvalidCredentials, commonkeys.Username, usernameReq)
		return userDomain.User{}, "", "", errors.New(UserNotFoundOrInvalidCredentials)
	}

	span.AddEvent(EventComparePassword)
	if err := s.hasher.Compare(user.Password, passwordReq); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, InvalidCredentials)
		s.logger.WarnwCtx(ctx, ErrorToCompareHashAndPassword, commonkeys.Username, user.Username)
		return userDomain.User{}, "", "", errors.New(InvalidCredentials)
	}

	roles := user.Roles
	claims := map[string]any{
		claimskeys.Roles: roles,
		claimskeys.Name:  user.Name,
	}

	// Delegate token generation and storage to reduce function length.
	accessToken, refreshToken, err := s.generateAndStoreTokens(ctx, user.ID, claims)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToCreateToken)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error())
		return userDomain.User{}, "", "", errors.New(ErrorToCreateToken)
	}

	span.AddEvent(EventLoginSuccess)
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))
	span.SetStatus(codes.Ok, SuccessToLogin)

	s.logger.InfowCtx(ctx, SuccessToLogin, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, accessToken, refreshToken, nil
}

// helper to compute TTL from claims 'exp' value.
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

// generateAndStoreTokens generates access and refresh tokens, saves them in the auth store
// and returns the token values. This helper keeps `Login` concise while preserving
// observability and error handling.
func (s *Service) generateAndStoreTokens(ctx context.Context, userID uint64, claims map[string]any) (string, string, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, "generateAndStoreTokens")
	defer span.End()

	span.AddEvent(EventGenerateToken)
	accessValue, err := s.authProvider.GenerateAccessToken(userID, claims)
	if err != nil {
		return "", "", err
	}

	// compute TTL for access token by parsing claims via Verify
	accessTTL := time.Duration(0)
	if accessValue != "" {
		if c, err := s.authProvider.Verify(accessValue); err == nil {
			accessTTL = claimsTTLFromVerifyResult(c)
		}
	}

	accessAuth := authDomain.NewAccessToken(accessValue, userID).ToAuth()
	span.AddEvent(EventSaveTokenToStore)
	if accessTTL > 0 {
		if err := s.authStore.Save(ctx, accessAuth, accessTTL); err != nil {
			return "", "", err
		}
	}

	span.AddEvent("generate_refresh_token")
	refreshValue, err := s.authProvider.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	// compute TTL for refresh token
	refreshTTL := time.Duration(0)
	if refreshValue != "" {
		if c, err := s.authProvider.Verify(refreshValue); err == nil {
			refreshTTL = claimsTTLFromVerifyResult(c)
		}
	}

	refreshAuth := authDomain.NewRefreshToken(refreshValue, userID).ToAuth()
	if refreshTTL > 0 {
		if err := s.authStore.Save(ctx, refreshAuth, refreshTTL); err != nil {
			return "", "", err
		}
	}

	return accessValue, refreshValue, nil
}
