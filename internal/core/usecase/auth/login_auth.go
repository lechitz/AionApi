// internal/core/usecase/auth/login.go
// Package auth contains use cases for authenticating users and generating tokens.
package auth

import (
	"context"
	"errors"
	"strconv"

	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/core/usecase/auth/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
func (s *Service) Login(ctx context.Context, usernameReq, passwordReq string) (domain.User, string, error) {
	tracer := otel.Tracer(constants.TracerName)
	ctx, span := tracer.Start(ctx, constants.SpanLogin)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, constants.SpanLogin),
		attribute.String(commonkeys.Username, usernameReq),
	)

	span.AddEvent(constants.EventLookupUser)
	user, err := s.userRepository.GetByUsername(ctx, usernameReq)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToGetUserByUserName)
		s.logger.ErrorwCtx(ctx, constants.ErrorToGetUserByUserName, commonkeys.Error, err.Error())
		return domain.User{}, "", errors.New(constants.ErrorToGetUserByUserName)
	}
	if user.ID == 0 {
		span.SetStatus(codes.Error, constants.UserNotFoundOrInvalidCredentials)
		s.logger.WarnwCtx(ctx, constants.UserNotFoundOrInvalidCredentials, commonkeys.Username, usernameReq)
		return domain.User{}, "", errors.New(constants.UserNotFoundOrInvalidCredentials)
	}

	span.AddEvent(constants.EventComparePassword)
	if err := s.hasher.Compare(user.Password, passwordReq); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.InvalidCredentials)
		s.logger.WarnwCtx(ctx, constants.ErrorToCompareHashAndPassword, commonkeys.Username, user.Username)
		return domain.User{}, "", errors.New(constants.InvalidCredentials)
	}

	span.AddEvent(constants.EventGenerateToken)
	tokenValue, err := s.tokenProvider.Generate(ctx, user.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToCreateToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error())
		return domain.User{}, "", errors.New(constants.ErrorToCreateToken)
	}

	span.AddEvent(constants.EventSaveTokenToStore)
	token := domain.Token{Key: user.ID, Value: tokenValue}
	if err := s.tokenStore.Save(ctx, token); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, constants.ErrorToCreateToken)
		s.logger.ErrorwCtx(ctx, constants.ErrorToCreateToken, commonkeys.Error, err.Error())
		return domain.User{}, "", errors.New(constants.ErrorToCreateToken)
	}

	span.AddEvent(constants.EventLoginSuccess)
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))
	span.SetStatus(codes.Ok, constants.SuccessToLogin)

	s.logger.InfowCtx(ctx, constants.SuccessToLogin, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, tokenValue, nil
}
