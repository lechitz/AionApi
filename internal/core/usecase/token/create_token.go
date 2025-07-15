// Package token contains use cases for managing tokens in the system.
package token

import (
	"context"
	"github.com/lechitz/AionApi/internal/adapters/secondary/security/jwt"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/usecase/token/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// CreateToken generates and persists a new signed token for the user, replacing any existing token.
func (s *Service) CreateToken(ctx context.Context, userID uint64) (string, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanCreateToken)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanCreateToken),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	token, err := s.tokenRepository.Get(ctx, userID)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToGetToken))
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetToken, commonkeys.Error, err.Error())
		return "", err
	}

	if token != "" {
		if err := s.tokenRepository.Delete(ctx, token); err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToDeleteToken))
			s.logger.ErrorwCtx(ctx, constants.ErrorToDeleteToken, commonkeys.Error, err.Error())
			return "", err
		}
	}

	signedToken, err := jwt.GenerateToken(tokenDomain.UserID, s.secretKey)
	if err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToAssignToken))
		s.logger.ErrorwCtx(ctx, constants.ErrorToAssignToken, commonkeys.Error, err.Error())
		return "", err
	}

	tokenDomain.Token = signedToken

	if err := s.tokenRepository.Save(ctx, tokenDomain); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String(commonkeys.Status, constants.ErrorToSaveToken))
		s.logger.ErrorwCtx(ctx, constants.ErrorToSaveToken, commonkeys.Error, err.Error())
		return "", err
	}

	span.SetAttributes(
		attribute.String(commonkeys.Status, constants.SuccessTokenCreated),
		attribute.String(commonkeys.UserID, strconv.FormatUint(tokenDomain.UserID, 10)),
	)
	s.logger.InfowCtx(ctx, constants.SuccessTokenCreated, commonkeys.UserID, strconv.FormatUint(tokenDomain.UserID, 10))
	return signedToken, nil
}
