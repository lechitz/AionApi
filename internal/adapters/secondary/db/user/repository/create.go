package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/mapper"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/domain"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/shared/sharederrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create inserts a new user. Returns a ValidationError if username/email is already in use.
func (up UserRepository) Create(ctx context.Context, userDomain domain.User) (domain.User, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanCreate, trace.WithAttributes(
		attribute.String(commonkeys.Username, userDomain.Username),
		attribute.String(commonkeys.Email, userDomain.Email),
		attribute.String(commonkeys.Operation, constants.OperationCreate),
	))
	defer span.End()

	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).Create(&userDB).Error; err != nil {
		if field, ok := isUniqueViolation(err); ok {
			span.SetStatus(codes.Error, constants.StatusValidationDuplicate)
			span.SetAttributes(attribute.String(constants.AttrHTTPErrorReason, field+constants.SuffixAlreadyExists))
			span.RecordError(err)
			up.logger.ErrorwCtx(ctx, constants.LogUniqueViolationOnCreate, constants.LogField, field, commonkeys.Error, err.Error())
			return domain.User{}, sharederrors.NewValidationError(field, field+constants.MsgAlreadyInUse)
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, constants.LogFailedCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, constants.StatusUserCreated)
	up.logger.InfowCtx(ctx, constants.LogUserCreated, commonkeys.UserID, userDB.ID)
	return mapper.UserFromDB(userDB), nil
}
