package repository

import (
	"context"
	"errors"
	"time"

	"github.com/lechitz/aion-api/internal/shared/constants/commonkeys"
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/mapper"
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	userOutput "github.com/lechitz/aion-api/internal/user/core/ports/output"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

const (
	registrationPendingStatus = "pending"
)

// CreateRegistrationSession persists a staged public registration session.
func (up UserRepository) CreateRegistrationSession(ctx context.Context, session domain.RegistrationSession) (domain.RegistrationSession, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.registration_session.create")
	defer span.End()

	sessionDB := mapper.RegistrationSessionToDB(session)
	if err := up.db.WithContext(ctx).Create(&sessionDB).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.RegistrationSession{}, err
	}

	span.SetStatus(codes.Ok, "registration_session_created")
	return mapper.RegistrationSessionFromDB(sessionDB), nil
}

// GetRegistrationSessionByID retrieves a staged registration session by id.
func (up UserRepository) GetRegistrationSessionByID(ctx context.Context, registrationID string) (domain.RegistrationSession, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.registration_session.get_by_id")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

	var sessionDB model.RegistrationSessionDB
	err := up.db.WithContext(ctx).
		Where("registration_id = ?", registrationID).
		First(&sessionDB).Error()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			span.SetStatus(codes.Ok, "registration_session_not_found")
			return domain.RegistrationSession{}, nil
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.RegistrationSession{}, err
	}

	span.SetStatus(codes.Ok, "registration_session_retrieved")
	return mapper.RegistrationSessionFromDB(sessionDB), nil
}

// UpdateRegistrationSession updates selected fields and returns the updated staged session.
func (up UserRepository) UpdateRegistrationSession(ctx context.Context, registrationID string, fields map[string]interface{}) (domain.RegistrationSession, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.registration_session.update")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

	if err := up.db.WithContext(ctx).
		Model(&model.RegistrationSessionDB{}).
		Where("registration_id = ?", registrationID).
		Updates(fields).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return domain.RegistrationSession{}, err
	}

	return up.GetRegistrationSessionByID(ctx, registrationID)
}

// DeleteRegistrationSession removes the staged registration session.
func (up UserRepository) DeleteRegistrationSession(ctx context.Context, registrationID string) error {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.registration_session.delete")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

	if err := up.db.WithContext(ctx).
		Where("registration_id = ?", registrationID).
		Delete(&model.RegistrationSessionDB{}).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	span.SetStatus(codes.Ok, "registration_session_deleted")
	return nil
}

// CheckRegistrationUniqueness checks if username/email already exists in pending non-expired sessions.
func (up UserRepository) CheckRegistrationUniqueness(ctx context.Context, username, email string, now time.Time) (userOutput.RegistrationSessionUniqueness, error) {
	tr := otel.Tracer(TracerName)
	ctx, span := tr.Start(ctx, "user.registration_session.check_uniqueness")
	defer span.End()
	span.SetAttributes(
		attribute.String(commonkeys.Username, username),
		attribute.String(commonkeys.Email, email),
	)

	var usernameCount int64
	if err := up.db.WithContext(ctx).
		Model(&model.RegistrationSessionDB{}).
		Where("username = ? AND status = ? AND expires_at > ?", username, registrationPendingStatus, now).
		Count(&usernameCount).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return userOutput.RegistrationSessionUniqueness{}, err
	}

	var emailCount int64
	if err := up.db.WithContext(ctx).
		Model(&model.RegistrationSessionDB{}).
		Where("email = ? AND status = ? AND expires_at > ?", email, registrationPendingStatus, now).
		Count(&emailCount).Error(); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return userOutput.RegistrationSessionUniqueness{}, err
	}

	span.SetStatus(codes.Ok, "registration_uniqueness_checked")
	return userOutput.RegistrationSessionUniqueness{
		UsernameTaken: usernameCount > 0,
		EmailTaken:    emailCount > 0,
	}, nil
}
