// Package token contains use cases for managing tokens in the system.
package token

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create generates and persists a new signed token for the user, replacing any existing token.
func (s *Service) Create(ctx context.Context, tokenKey uint64) (domain.Token, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(
		ctx,
		constants.SpanCreateToken,
		trace.WithAttributes(
			attribute.String(commonkeys.Operation, constants.SpanCreateToken),
			attribute.String(commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10)),
		),
	)
	defer span.End()

	span.AddEvent(constants.EventGenerateToken)
	newTokenValue, err := s.tokenProvider.Generate(ctx, tokenKey)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToAssignToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToAssignToken, commonkeys.Error, err.Error())
		return domain.Token{}, fmt.Errorf("%s: %w", constants.ErrorToAssignToken, err)
	}
	if newTokenValue == "" {
		err := errors.New(constants.ErrEmptyTokenGenerated)
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrEmptyTokenGenerated)
		s.logger.ErrorwCtx(ctx, constants.ErrEmptyTokenGenerated,
			commonkeys.TokenKey, strconv.FormatUint(tokenKey, 10),
		)
		return domain.Token{}, err
	}

	token := domain.Token{Key: tokenKey, Value: newTokenValue}

	span.AddEvent(constants.EventSaveTokenToStore)
	if err := s.tokenStore.Save(ctx, token); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToSaveToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToSaveToken, commonkeys.Error, err.Error())
		return domain.Token{}, fmt.Errorf("%s: %w", constants.ErrorToSaveToken, err)
	}

	span.AddEvent(constants.EventTokenCreated)
	span.SetStatus(codes.Ok, constants.SuccessTokenCreated)

	s.logger.InfowCtx(ctx, constants.SuccessTokenCreated,
		commonkeys.TokenKey, strconv.FormatUint(token.Key, 10),
	)
	return token, nil
}
