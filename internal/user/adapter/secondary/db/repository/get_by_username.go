package repository

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GetByUsername retrieves a user by username. Returns a zero-value user if not found. //nolint:dupl.
//
//nolint:dupl
func (up UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	tr := otel.Tracer(TracerUserRepository)
	ctx, span := tr.Start(ctx, SpanGetByUsername, trace.WithAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String(commonkeys.Operation, OperationGetByUsername),
	))
	defer span.End()

	var userDB model.UserDB
	err := up.db.WithContext(ctx).
		Select(SelectByUsernameColumns).
		Where(commonkeys.Username+" = ?", username).
		First(&userDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, StatusUserNotFoundOK)
			up.logger.InfowCtx(ctx, LogUserNotFoundByUsername, commonkeys.Username, username)
			return domain.User{}, nil
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, LogFailedGetByUsername, commonkeys.Username, username, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, StatusUserRetrievedByUsername)
	up.logger.InfowCtx(ctx, LogUserRetrievedByUsername, commonkeys.UserID, userDB.ID, commonkeys.Username, userDB.Username)
	return mapper.UserFromDB(userDB), nil
}
