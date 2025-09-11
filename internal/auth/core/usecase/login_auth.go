// Package usecase (auth) contains use cases for authenticating users and generating tokens.
package usecase

import (
	"context"
	"errors"
	"strconv"

	authDomain "github.com/lechitz/AionApi/internal/auth/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	userDomain "github.com/lechitz/AionApi/internal/user/core/domain"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Login authenticates a user by validating credentials and generates a new token if valid.
func (s *Service) Login(ctx context.Context, usernameReq, passwordReq string) (userDomain.User, string, error) {
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
		return userDomain.User{}, "", errors.New(ErrorToGetUserByUserName)
	}
	if user.ID == 0 {
		span.SetStatus(codes.Error, UserNotFoundOrInvalidCredentials)
		s.logger.WarnwCtx(ctx, UserNotFoundOrInvalidCredentials, commonkeys.Username, usernameReq)
		return userDomain.User{}, "", errors.New(UserNotFoundOrInvalidCredentials)
	}

	span.AddEvent(EventComparePassword)
	if err := s.hasher.Compare(user.Password, passwordReq); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, InvalidCredentials)
		s.logger.WarnwCtx(ctx, ErrorToCompareHashAndPassword, commonkeys.Username, user.Username)
		return userDomain.User{}, "", errors.New(InvalidCredentials)
	}

	span.AddEvent(EventGenerateToken)
	tokenValue, err := s.authProvider.Generate(ctx, user.ID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToCreateToken)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error())
		return userDomain.User{}, "", errors.New(ErrorToCreateToken)
	}

	token := authDomain.Auth{
		Key:   user.ID,
		Token: tokenValue,
	}

	span.AddEvent(EventSaveTokenToStore)
	if err := s.authStore.Save(ctx, token); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, ErrorToCreateToken)
		s.logger.ErrorwCtx(ctx, ErrorToCreateToken, commonkeys.Error, err.Error())
		return userDomain.User{}, "", errors.New(ErrorToCreateToken)
	}

	span.AddEvent(EventLoginSuccess)
	span.SetAttributes(attribute.String(commonkeys.UserID, strconv.FormatUint(user.ID, 10)))
	span.SetStatus(codes.Ok, SuccessToLogin)

	s.logger.InfowCtx(ctx, SuccessToLogin, commonkeys.UserID, strconv.FormatUint(user.ID, 10))
	return user, tokenValue, nil
}
