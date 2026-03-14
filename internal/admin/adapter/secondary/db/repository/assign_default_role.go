// Package repository provides admin database operations.
package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/lechitz/AionApi/internal/platform/ports/output/db"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AssignDefaultRole assigns the default 'user' role to a newly created user.
// This is called by other bounded contexts (e.g., /user) during user creation.
func (r *AdminRepository) AssignDefaultRole(ctx context.Context, userID uint64) error {
	return r.assignDefaultRole(ctx, r.db, userID)
}

// AssignDefaultRoleWithTx assigns the default role using the provided transaction handle.
func (r *AdminRepository) AssignDefaultRoleWithTx(ctx context.Context, tx db.DB, userID uint64) error {
	if tx == nil {
		tx = r.db
	}
	return r.assignDefaultRole(ctx, tx, userID)
}

func (r *AdminRepository) assignDefaultRole(ctx context.Context, database db.DB, userID uint64) error {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, SpanAssignDefaultRole)
	defer span.End()

	span.SetAttributes(
		attribute.String(commonkeys.Operation, OperationAssignDefaultRole),
		attribute.String(commonkeys.UserID, strconv.FormatUint(userID, 10)),
	)

	r.logger.InfowCtx(ctx, LogAssigningDefaultRole, commonkeys.UserID, userID)

	// Check if role exists
	var count int64
	if err := database.WithContext(ctx).Raw(QueryCountDefaultRole, DefaultRoleName).Scan(&count).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, LogFailedToAssignDefaultRole, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return fmt.Errorf("%w: %w", ErrCheckDefaultRoleExists, err)
	}

	if count == 0 {
		err := ErrDefaultRoleNotFound
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, LogFailedToAssignDefaultRole, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return err
	}

	// Insert user-role assignment
	now := time.Now().UTC()
	if err := database.WithContext(ctx).Exec(QueryInsertDefaultRoleAssignment, userID, now, DefaultRoleName).Error(); err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		r.logger.ErrorwCtx(ctx, LogFailedToAssignDefaultRole, commonkeys.UserID, userID, commonkeys.Error, err.Error())
		return fmt.Errorf("%w: %w", ErrAssignDefaultRole, err)
	}

	span.SetStatus(codes.Ok, StatusDefaultRoleAssigned)
	r.logger.InfowCtx(ctx, LogDefaultRoleAssigned, commonkeys.UserID, userID)

	return nil
}
