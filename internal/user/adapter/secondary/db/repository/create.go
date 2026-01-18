package repository

import (
	"context"
	"fmt"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type roleAssignerWithTx interface {
	AssignDefaultRoleWithTx(ctx context.Context, tx db.DB, userID uint64) error
}

// Create inserts a new user and assigns default roles. Returns a ValidationError if username/email is already in use.
func (up UserRepository) Create(ctx context.Context, userDomain domain.User) (domain.User, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, SpanCreate, trace.WithAttributes(
		attribute.String(commonkeys.Username, userDomain.Username),
		attribute.String(commonkeys.Email, userDomain.Email),
		attribute.String(commonkeys.Operation, OperationCreate),
	))
	defer span.End()

	userDB := mapper.UserToDB(userDomain)

	// Use transaction to ensure atomicity
	err := up.db.WithContext(ctx).Transaction(func(tx db.DB) error {
		// Create user
		if err := tx.Create(&userDB).Error(); err != nil {
			return err
		}

		// Delegate default role assignment to /admin context
		// This maintains separation of concerns: /user creates users, /admin manages roles
		if assigner, ok := up.roleAssigner.(roleAssignerWithTx); ok {
			if err := assigner.AssignDefaultRoleWithTx(ctx, tx, userDB.ID); err != nil {
				return fmt.Errorf("failed to assign default role: %w", err)
			}
		} else if err := up.roleAssigner.AssignDefaultRole(ctx, userDB.ID); err != nil {
			return fmt.Errorf("failed to assign default role: %w", err)
		}

		return nil
	})
	if err != nil {
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

	// Fetch the created user
	createdUser, err := up.GetByID(ctx, userDB.ID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		return domain.User{}, err
	}

	span.SetStatus(codes.Ok, StatusUserCreated)
	up.logger.InfowCtx(ctx, LogUserCreated, commonkeys.UserID, userDB.ID)
	return createdUser, nil
}
