package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/user/model"
	constants "github.com/lechitz/AionApi/internal/adapters/secondary/db/user/repository/constants"
	"github.com/lechitz/AionApi/internal/core/ports/output"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// CheckUniqueness verifies whether username and/or email are already taken.
// It returns which fields are taken and (optionally) the owner IDs.
func (up UserRepository) CheckUniqueness(ctx context.Context, username, email string) (output.UserUniqueness, error) {
	tr := otel.Tracer(constants.TracerUserRepository)
	ctx, span := tr.Start(ctx, constants.SpanCheckUniqueness, trace.WithAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String(commonkeys.Email, email),
		attribute.String(commonkeys.Operation, constants.OperationCheckUniqueness),
	))
	defer span.End()

	var res output.UserUniqueness

	lookupIDByField := func(field, value, logOnFail string) (*uint64, error) {
		if strings.TrimSpace(value) == "" {
			return nil, nil
		}

		var u model.UserDB
		err := up.db.WithContext(ctx).
			Model(&model.UserDB{}).
			Select(commonkeys.UserID).
			Where(field+" = ?", value).
			First(&u).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			up.logger.ErrorwCtx(ctx, logOnFail, field, value, commonkeys.Error, err.Error())
			return nil, err
		}

		id := u.ID
		return &id, nil
	}

	if id, err := lookupIDByField(commonkeys.Username, username, constants.LogFailedCheckUsername); err != nil {
		return res, err
	} else if id != nil {
		res.UsernameTaken = true
		res.UsernameOwnerID = id
	}

	if id, err := lookupIDByField(commonkeys.Email, email, constants.LogFailedCheckEmail); err != nil {
		return res, err
	} else if id != nil {
		res.EmailTaken = true
		res.EmailOwnerID = id
	}

	span.SetStatus(codes.Ok, constants.StatusUniquenessChecked)
	return res, nil
}
