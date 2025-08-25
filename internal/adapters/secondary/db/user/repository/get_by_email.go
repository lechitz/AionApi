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

// GetByEmail retrieves a user by email. Returns zero-value user if not found.
func (up UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanGetByEmail, trace.WithAttributes(
		attribute.String(commonkeys.Email, email),
		attribute.String(commonkeys.Operation, constants.OperationGetByEmail),
	))
	defer span.End()

	var userDB model.UserDB
	err := up.db.WithContext(ctx).
		Select(constants.SelectByEmailColumns).
		Where(commonkeys.Email+" = ?", email).
		First(&userDB).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, constants.StatusUserNotFoundOK)
			up.logger.InfowCtx(ctx, constants.LogUserNotFoundByEmail, commonkeys.Email, email)
			return domain.User{}, nil
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedGetByEmail, commonkeys.Email, email, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, constants.StatusUserRetrievedByEmail)
	up.logger.InfowCtx(ctx, constants.LogUserRetrievedByEmail, commonkeys.UserID, userDB.ID, commonkeys.Email, userDB.Email)
	return mapper.UserFromDB(userDB), nil
}
