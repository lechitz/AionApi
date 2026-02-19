package usecase

import (
	"context"
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

func (s *Service) StartRegistration(ctx context.Context, cmd input.StartRegistrationCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, "user.registration.start")
	defer span.End()

	if s.registrationRepo == nil {
		return domain.RegistrationSession{}, fmt.Errorf("registration repository not configured")
	}

	name := strings.TrimSpace(cmd.Name)
	username := strings.ToLower(strings.TrimSpace(cmd.Username))
	email := strings.ToLower(strings.TrimSpace(cmd.Email))
	password := cmd.Password

	if name == "" || username == "" || email == "" || password == "" {
		return domain.RegistrationSession{}, sharederrors.MissingFields(commonkeys.Name, commonkeys.Username, commonkeys.Email, commonkeys.Password)
	}
	if len(password) < 8 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Password, "password must be at least 8 characters")
	}
	if !emailRegex.MatchString(email) {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Email, "invalid email format")
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
		CurrentStep:    1,
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
	span.SetAttributes(attribute.String("registration_id", out.RegistrationID))
	span.SetStatus(codes.Ok, "registration_started")
	return out, nil
}

func (s *Service) UpdateRegistrationProfile(ctx context.Context, registrationID string, cmd input.UpdateRegistrationProfileCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, "user.registration.update_profile")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

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
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Locale, "invalid locale format")
	}
	if len(timezone) > 64 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Timezone, "timezone must be up to 64 characters")
	}
	if len(location) > 255 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Location, "location must be up to 255 characters")
	}
	if len(bio) > 1000 {
		return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.Bio, "bio must be up to 1000 characters")
	}

	if session.CurrentStep < 1 {
		return domain.RegistrationSession{}, fmt.Errorf("%w: invalid registration step", sharederrors.ErrDomainConflict)
	}

	fields := map[string]interface{}{
		commonkeys.Locale:   locale,
		commonkeys.Timezone: timezone,
		commonkeys.Location: location,
		"current_step":      2,
		"updated_at":        time.Now().UTC(),
	}
	if bio != "" {
		fields[commonkeys.Bio] = bio
	}
	out, err := s.registrationRepo.UpdateRegistrationSession(ctx, registrationID, fields)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	span.SetStatus(codes.Ok, "registration_profile_updated")
	return out, nil
}

func (s *Service) UpdateRegistrationAvatar(ctx context.Context, registrationID string, cmd input.UpdateRegistrationAvatarCommand) (domain.RegistrationSession, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, "user.registration.update_avatar")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

	session, err := s.getActiveRegistrationSession(ctx, registrationID)
	if err != nil {
		span.RecordError(err)
		return domain.RegistrationSession{}, err
	}
	if session.CurrentStep < 2 {
		return domain.RegistrationSession{}, fmt.Errorf("%w: profile step must be completed before avatar step", sharederrors.ErrDomainConflict)
	}

	fields := map[string]interface{}{
		"current_step": 3,
		"updated_at":   time.Now().UTC(),
	}

	if cmd.AvatarURL != nil {
		trimmed := strings.TrimSpace(*cmd.AvatarURL)
		if trimmed != "" {
			if _, parseErr := url.ParseRequestURI(trimmed); parseErr != nil {
				return domain.RegistrationSession{}, sharederrors.NewValidationError(commonkeys.AvatarURL, "avatar_url must be a valid URL")
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
	span.SetStatus(codes.Ok, "registration_avatar_updated")
	return out, nil
}

func (s *Service) CompleteRegistration(ctx context.Context, registrationID string) (domain.User, error) {
	tracer := otel.Tracer(TracerName)
	ctx, span := tracer.Start(ctx, "user.registration.complete")
	defer span.End()
	span.SetAttributes(attribute.String("registration_id", registrationID))

	session, err := s.getActiveRegistrationSession(ctx, registrationID)
	if err != nil {
		span.RecordError(err)
		return domain.User{}, err
	}
	if session.CurrentStep < 3 {
		return domain.User{}, fmt.Errorf("%w: registration flow not completed", sharederrors.ErrDomainConflict)
	}
	if session.Locale == nil || session.Timezone == nil || session.Location == nil {
		return domain.User{}, fmt.Errorf("%w: required profile fields are missing", sharederrors.ErrDomainConflict)
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
		s.logger.WarnwCtx(ctx, "failed to delete completed registration session", "registration_id", registrationID, commonkeys.Error, err)
	}

	span.SetStatus(codes.Ok, "registration_completed")
	return user, nil
}

func (s *Service) getActiveRegistrationSession(ctx context.Context, registrationID string) (domain.RegistrationSession, error) {
	if s.registrationRepo == nil {
		return domain.RegistrationSession{}, fmt.Errorf("registration repository not configured")
	}
	if strings.TrimSpace(registrationID) == "" {
		return domain.RegistrationSession{}, sharederrors.NewValidationError("registration_id", "registration_id is required")
	}
	session, err := s.registrationRepo.GetRegistrationSessionByID(ctx, registrationID)
	if err != nil {
		return domain.RegistrationSession{}, err
	}
	if session.RegistrationID == "" {
		return domain.RegistrationSession{}, sharederrors.NewValidationError("registration_id", "registration session not found")
	}
	if session.Status != domain.RegistrationStatusPending {
		return domain.RegistrationSession{}, fmt.Errorf("%w: registration session is not pending", sharederrors.ErrDomainConflict)
	}
	if time.Now().UTC().After(session.ExpiresAt) {
		return domain.RegistrationSession{}, fmt.Errorf("%w: registration session expired", sharederrors.ErrDomainConflict)
	}
	return session, nil
}
