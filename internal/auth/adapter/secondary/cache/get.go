package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Get returns the token for a user or empty if not present.
func (s *Store) Get(ctx context.Context, tokenKey uint64, tokenType string) (domain.Auth, error) {
	_ = tokenType // explicit usage to satisfy static analyzers when tokenType may be conditionally used
	tr := otel.Tracer(SpanTracerTokenStore)
	ctx, span := tr.Start(ctx, SpanNameTokenGet, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(tokenKey, 10)),
		attribute.String(commonkeys.Operation, OperationGet),
		attribute.String(commonkeys.Entity, commonkeys.EntityToken),
	))
	defer span.End()

	cacheKey := fmt.Sprintf(TokenUserKeyFormat, tokenKey, tokenType)

	tokenValue, err := s.cache.Get(ctx, cacheKey)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		s.logger.Errorw(ErrorToGetTokenFromRedis, commonkeys.TokenKey, cacheKey, commonkeys.Error, err)
		return domain.Auth{}, err
	}

	if tokenValue == "" {
		span.SetStatus(codes.Ok, sharederrors.ErrTokenNotFound)
		return domain.Auth{}, nil
	}

	var tokenDomain domain.Auth
	switch tokenType {
	case commonkeys.TokenTypeAccess:
		tokenDomain = domain.NewAccessToken(tokenValue, tokenKey).ToAuth()
	case commonkeys.TokenTypeRefresh:
		tokenDomain = domain.NewRefreshToken(tokenValue, tokenKey).ToAuth()
	default:
		// Fallback to generic Auth when tokenType is unknown
		tokenDomain = domain.Auth{Key: tokenKey, Token: tokenValue, Type: tokenType}
	}

	span.SetStatus(codes.Ok, TokenRetrievedSuccessfully)
	return tokenDomain, nil
}
