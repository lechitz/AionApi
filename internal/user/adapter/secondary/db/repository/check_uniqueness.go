package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// CheckUniqueness verifies whether username and/or email are already taken.
// It returns which fields are taken and (optionally) the owner IDs.
func (up UserRepository) CheckUniqueness(ctx context.Context, username, email string) (output.UserUniqueness, error) {
	tr := otel.Tracer(TracerUserRepository)
	ctx, span := tr.Start(ctx, SpanCheckUniqueness, trace.WithAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String(commonkeys.Email, email),
		attribute.String(commonkeys.Operation, OperationCheckUniqueness),
	))
	defer span.End()

	var res output.UserUniqueness

	// Returns (id, found, err). Avoid it (nil, nil) by not using pointer in the helper.
	lookupIDByField := func(field, value, logOnFail string) (uint64, bool, error) {
		if strings.TrimSpace(value) == "" {
			return 0, false, nil
		}

		var u model.UserDB
		err := up.db.WithContext(ctx).
			Model(&model.UserDB{}).
			Select(commonkeys.UserID).
			Where(field+" = ?", value).
			First(&u).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, false, nil
		}
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
			span.RecordError(err)
			up.logger.ErrorwCtx(ctx, logOnFail, field, value, commonkeys.Error, err.Error())
			return 0, false, err
		}

		return u.ID, true, nil
	}

	if id, found, err := lookupIDByField(commonkeys.Username, username, LogFailedCheckUsername); err != nil {
		return res, err
	} else if found {
		res.UsernameTaken = true
		res.UsernameOwnerID = &id
	}

	if id, found, err := lookupIDByField(commonkeys.Email, email, LogFailedCheckEmail); err != nil {
		return res, err
	} else if found {
		res.EmailTaken = true
		res.EmailOwnerID = &id
	}

	span.SetStatus(codes.Ok, StatusUniquenessChecked)
	return res, nil
}
