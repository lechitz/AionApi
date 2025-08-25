package repository

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/mapper"
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/model"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// GetByID retrieves a user by ID.
func (up UserRepository) GetByID(ctx context.Context, userID uint64) (domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanGetByID, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, constants.OperationGetByID),
	))
	defer span.End()

	var userDB model.UserDB
	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where(commonkeys.UserID+" = ?", userID).
		First(&userDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedGetByID, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, constants.StatusUserRetrievedByID)
	up.logger.InfowCtx(ctx, constants.LogUserRetrievedByID, commonkeys.UserID, userDB.ID)
	return mapper.UserFromDB(userDB), nil
}
