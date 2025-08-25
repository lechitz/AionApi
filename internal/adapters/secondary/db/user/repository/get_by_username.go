package repository

import (
	"context"
	"errors"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/model"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// GetByUsername retrieves a user by username. Returns zero-value user if not found.
func (up UserRepository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanGetByUsername, trace.WithAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String(commonkeys.Operation, constants.OperationGetByUsername),
	))
	defer span.End()

	var userDB model.UserDB
	err := up.db.WithContext(ctx).
		Select(constants.SelectByUsernameColumns).
		Where(commonkeys.Username+" = ?", username).
		First(&userDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, constants.StatusUserNotFoundOK)
			up.logger.InfowCtx(ctx, constants.LogUserNotFoundByUsername, commonkeys.Username, username)
			return domain.User{}, nil
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedGetByUsername, commonkeys.Username, username, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, constants.StatusUserRetrievedByUsername)
	up.logger.InfowCtx(ctx, constants.LogUserRetrievedByUsername, commonkeys.UserID, userDB.ID, commonkeys.Username, userDB.Username)
	return mapper.UserFromDB(userDB), nil
}
