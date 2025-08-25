package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/model"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SoftDelete marks a user as deleted by setting DeletedAt and UpdatedAt.
func (up UserRepository) SoftDelete(ctx context.Context, userID uint64) error {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanSoftDelete, trace.WithAttributes(
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
		attribute.String(commonkeys.Operation, constants.OperationSoftDelete),
	))
	defer span.End()

	fields := map[string]interface{}{
		commonkeys.UserDeletedAt: time.Now().UTC(),
		commonkeys.UserUpdatedAt: time.Now().UTC(),
	}

	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Where(commonkeys.UserID+" = ?", userID).
		Updates(fields).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedSoftDelete, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return err
	}

	span.SetStatus(codes.Ok, constants.StatusUserSoftDeleted)
	up.logger.InfowCtx(ctx, constants.LogUserSoftDeleted, commonkeys.UserID, userID)
	return nil
}
