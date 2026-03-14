package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lechitz/AionApi/internal/platform/server/http/utils/sharederrors"
	"github.com/lechitz/AionApi/internal/shared/constants/commonkeys"
	"github.com/lechitz/AionApi/internal/user/core/domain"
	"github.com/lechitz/AionApi/internal/user/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const (
	registrationSessionTTL = 2 * time.Hour
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// StartRegistration validates user identity data and starts a registration session.
func (s *Service) StartRegistration(ctx context.Context, cmd input.StartRegistrationCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, spanStartRegistration)
	defer span.End()

	if s.registrationRepo == nil {
		return domain.RegistrationSession{}, errors.New(errRegistrationRepoNotConfigured)
	}

	name := strings.TrimSpace(cmd.Name)
	username := strings.ToLower(strings.TrimSpace(cmd.Username))
	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	password := cmd.Password

	if name == "" || username == "" || email == "" || password == "" {
		return domain.RegistrationSession{}, sharederrors.MissingFields(commonkeys.Name, commonkeys.Username, commonkeys.Email, commonkeys.Password)
	}
	if len(password) < 8 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Password, errPasswordMinLength)
	}
	if !emailRegex.MatchString(email) {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Email, errInvalidEmailFormat)
	}

	conflict, err := s.userRepository.CheckUniqueness(ctx, username, email)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	if conflict.UsernameTaken {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
	}
	if conflict.EmailTaken {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
	}

	regConflict, err := s.registrationRepo.CheckRegistrationUniqueness(ctx, username, email, time.Now().UTC())
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	if regConflict.UsernameTaken {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
	}
	if regConflict.EmailTaken {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
	}

	hashed, err := s.hasher.Hash(password)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, ErrHashPassword
	}

	now := time.Now().UTC()
	session := domain.RegistrationSession{
		RegistrationID: uuid.NewString(),
		Name:           name,
		Username:       username,
		Email:          email,
		PasswordHash:   hashed,
		CurrentStep:    defaultRegistrationStep,
		Status:         domain.RegistrationStatusPending,
		ExpiresAt:      now.Add(registrationSessionTTL),
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	out, err := s.registrationRepo.CreateRegistrationSession(ctx, session)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	span.SetAttributes(attribute.String(attrRegistrationID, out.RegistrationID))
	span.SetStatus(codes.Ok, statusRegistrationStarted)
	return out, nil
}

// UpdateRegistrationProfile updates locale/timezone/location and optional bio for an active registration.
func (s *Service) UpdateRegistrationProfile(ctx context.Context, registrationID string, cmd input.UpdateRegistrationProfileCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, spanUpdateRegistrationProfile)
	defer span.End()
	span.SetAttributes(attribute.String(attrRegistrationID, registrationID))

	session, err := s.getActiveRegistrationSession(ctx, registrationID)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}

	locale := strings.TrimSpace(cmd.Locale)
	timezone := strings.TrimSpace(cmd.Timezone)
	location := strings.TrimSpace(cmd.Location)
	bio := strings.TrimSpace(cmd.Bio)

	if locale == "" || timezone == "" || location == "" {
		return domain.RegistrationSession{}, sharederrors.MissingFields(commonkeys.Locale, commonkeys.Timezone, commonkeys.Location)
	}
	if !regexp.MustCompile(`^[a-z]{2}(-[A-Z]{2})?$`).MatchString(locale) {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Locale, errInvalidLocaleFormat)
	}
	if len(timezone) > 64 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Timezone, errTimezoneTooLong)
	}
	if len(location) > 255 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Location, errLocationTooLong)
	}
	if len(bio) > 1000 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Bio, errBioTooLong)
	}

	if session.CurrentStep < defaultRegistrationStep {
		return domain.RegistrationSession{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errInvalidRegistrationStep)
	}

	fields := map[string]interface{}{
		commonkeys.Locale:   locale,
		commonkeys.Timezone: timezone,
		commonkeys.Location: location,
		attrCurrentStep:     profileRegistrationStep,
		attrUpdatedAt:       time.Now().UTC(),
	}
	if bio != "" {
		fields[commonkeys.Bio] = bio
	}
	out, err := s.registrationRepo.UpdateRegistrationSession(ctx, registrationID, fields)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	span.SetStatus(codes.Ok, statusRegistrationProfileUpdated)
	return out, nil
}

// UpdateRegistrationAvatar advances registration with an optional avatar URL.
func (s *Service) UpdateRegistrationAvatar(ctx context.Context, registrationID string, cmd input.UpdateRegistrationAvatarCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, spanUpdateRegistrationAvatar)
	defer span.End()
	span.SetAttributes(attribute.String(attrRegistrationID, registrationID))

	session, err := s.getActiveRegistrationSession(ctx, registrationID)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	if session.CurrentStep < profileRegistrationStep {
		return domain.RegistrationSession{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errProfileStepMustBeCompleted)
	}

	fields := map[string]interface{}{
		attrCurrentStep: avatarRegistrationStep,
		attrUpdatedAt:   time.Now().UTC(),
	}

	if cmd.AvatarURL != nil {
		trimmed := strings.TrimSpace(*cmd.AvatarURL)
		if trimmed != "" {
			if _, parseErr := url.ParseRequestURI(trimmed); parseErr != nil {
				return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.AvatarURL, errAvatarURLInvalid)
			}
			fields[commonkeys.AvatarURL] = trimmed
		} else {
			fields[commonkeys.AvatarURL] = nil
		}
	}

	out, err := s.registrationRepo.UpdateRegistrationSession(ctx, registrationID, fields)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	span.SetStatus(codes.Ok, statusRegistrationAvatarUpdated)
	return out, nil
}

// CompleteRegistration materializes the final user when all registration steps are complete.
func (s *Service) CompleteRegistration(ctx context.Context, registrationID string) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, spanCompleteRegistration)
	defer span.End()
	span.SetAttributes(attribute.String(attrRegistrationID, registrationID))

	session, err := s.getActiveRegistrationSession(ctx, registrationID)
	if err != nil {
		span.RecordError(err)
		return domain.User{}, err
	}
	if session.CurrentStep < avatarRegistrationStep {
		return domain.User{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errRegistrationFlowNotCompleted)
	}
	if session.Locale == nil || session.Timezone == nil || session.Location == nil {
		return domain.User{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errRegistrationRequiredFieldsMissing)
	}

	// Re-check uniqueness before materializing final user.
	conflict, err := s.userRepository.CheckUniqueness(ctx, session.Username, session.Email)
	if err != nil {
		span.RecordError(err)
		return domain.User{}, err
	}
	if conflict.UsernameTaken {
		return domain.User{}, sharederrors.NewValidationError(commonkeys.Username, sharederrors.ErrUsernameInUse)
	}
	if conflict.EmailTaken {
		return domain.User{}, sharederrors.NewValidationError(commonkeys.Email, sharederrors.ErrEmailInUse)
	}

	user, err := s.userRepository.Create(ctx, domain.User{
		Name:      session.Name,
		Username:  session.Username,
		Email:     session.Email,
		Password:  session.PasswordHash,
		Locale:    session.Locale,
		Timezone:  session.Timezone,
		Location:  session.Location,
		Bio:       session.Bio,
		AvatarURL: session.AvatarURL,
	})
	if err != nil {
		span.RecordError(err)
		return domain.User{}, err
	}

	if err := s.registrationRepo.DeleteRegistrationSession(ctx, registrationID); err != nil {
		s.logger.WarnwCtx(ctx, warnDeleteCompletedRegistrationSession, attrRegistrationID, registrationID, commonkeys.Error, err)
	}

	span.SetStatus(codes.Ok, statusRegistrationCompleted)
	return user, nil
}

func (s *Service) getActiveRegistrationSession(ctx context.Context, registrationID string) (domain.RegistrationSession, error) {
	if s.registrationRepo == nil {
		return domain.RegistrationSession{}, errors.New(errRegistrationRepoNotConfigured)
	}
	if strings.TrimSpace(registrationID) == "" {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(fieldRegistrationID, errRegistrationIDRequired)
	}
	session, err := s.registrationRepo.GetRegistrationSessionByID(ctx, registrationID)
	if err != nil {
		return domain.RegistrationSession{}, err
	}
	if session.RegistrationID == "" {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(fieldRegistrationID, errRegistrationSessionNotFound)
	}
	if session.Status != domain.RegistrationStatusPending {
		return domain.RegistrationSession{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errRegistrationSessionNotPending)
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return domain.RegistrationSession{}, fmt.Errorf("%w: %s", sharederrors.ErrDomainConflict, errRegistrationSessionExpired)
	}
	return session, nil
}
