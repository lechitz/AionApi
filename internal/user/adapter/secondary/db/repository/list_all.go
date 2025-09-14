package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ListAll returns all (non-deleted) users.
func (up UserRepository) ListAll(ctx context.Context) ([]domain.User, error) {
	tr := otel.Tracer(TracerUserRepository)
	ctx, span := tr.Start(ctx, SpanListAll, trace.WithAttributes(
		attribute.String(commonkeys.Operation, OperationListAll),
	))
	defer span.End()

	var usersDB []model.UserDB
	if err := up.db.WithContext(ctx).
		Model(&model.UserDB{}).
		Select(SelectListAllColumns).
		Find(&usersDB).Error; err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, LogFailedListAll, commonkeys.Error, err.Error())
		return nil, err
	}

	usersDomain := make([]domain.User, 0, len(usersDB))
	for _, u := range usersDB {
		usersDomain = append(usersDomain, mapper.UserFromDB(u))
	}

	span.SetStatus(codes.Ok, StatusUsersRetrieved)
	up.logger.InfowCtx(ctx, LogUsersRetrieved, commonkeys.UsersCount, len(usersDomain))
	return usersDomain, nil
}
