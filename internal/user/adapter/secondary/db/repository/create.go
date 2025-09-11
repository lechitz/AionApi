package repository

import (
	"context"

	"github.com/lechitz/AionApi/internal/platform/server/http/helpers/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Create inserts a new user. Returns a ValidationError if username/email is already in use.
func (up UserRepository) Create(ctx context.Context, userDomain domain.User) (domain.User, error) {
	tr := otel.Tracer(TracerUserRepository)
	ctx, span := tr.Start(ctx, SpanCreate, trace.WithAttributes(
		attribute.String(commonkeys.Username, userDomain.Username),
		attribute.String(commonkeys.Email, userDomain.Email),
		attribute.String(commonkeys.Operation, OperationCreate),
	))
	defer span.End()

	userDB := mapper.UserToDB(userDomain)

	if err := up.db.WithContext(ctx).Create(&userDB).Error; err != nil {
		if field, ok := isUniqueViolation(err); ok {
			span.SetStatus(codes.Error, StatusValidationDuplicate)
			span.SetAttributes(attribute.String(AttrHTTPErrorReason, field+SuffixAlreadyExists))
			span.RecordError(err)
			up.logger.ErrorwCtx(ctx, LogUniqueViolationOnCreate, LogField, field, commonkeys.Error, err.Error())
			return domain.User{}, sharederrors.NewValidationError(field, field+MsgAlreadyInUse)
		}
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		up.logger.ErrorwCtx(ctx, LogFailedCreateUser, commonkeys.Error, err.Error())
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, StatusUserCreated)
	up.logger.InfowCtx(ctx, LogUserCreated, commonkeys.UserID, userDB.ID)
	return mapper.UserFromDB(userDB), nil
}
