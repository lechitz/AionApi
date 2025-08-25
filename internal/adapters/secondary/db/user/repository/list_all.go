package repository

import (
	"context"

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

// ListAll returns all (non-deleted) users.
func (up UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanListAll, trace.WithAttributes(
		attribute.String(commonkeys.Operation, constants.OperationListAll),
	))
	defer span.End()

	var usersDB []model.UserDB
	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select(constants.SelectListAllColumns).
		Find(&usersDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedListAll, commonkeys.Error, err.Error())
		return nil, err
	}

	usersDomain := make([]domain.User, 0, len(usersDB))
	for _, u := range usersDB {
		usersDomain = append(usersDomain, mapper.UserFromDB(u))
	}

	span.SetStatus(codes.Ok, constants.StatusUsersRetrieved)
	up.logger.InfowCtx(ctx, constants.LogUsersRetrieved, commonkeys.UsersCount, len(usersDomain))
	return usersDomain, nil
}
