package repository

import (
	"context"
	"strconv"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/model"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Update updates fields for a user and returns the updated entity.
func (up UserRepository) Update(ctx context.Context, userID uint64, fields map[string]interface{}) (domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanUpdate, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, constants.OperationUpdate),
	))
	defer span.End()

	delete(fields, commonkeys.UserCreatedAt)

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where(commonkeys.UserID+" = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedUpdateUser, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, constants.StatusUserUpdated)
	up.logger.InfowCtx(ctx, constants.LogUserUpdated, commonkeys.UserID, userID)
	return up.GetByID(ctx, userID)
}
