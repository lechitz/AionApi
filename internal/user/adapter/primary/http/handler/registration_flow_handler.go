package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lechitz/aion-api/internal/platform/server/http/utils/httpresponse"
	"github.com/lechitz/aion-api/internal/user/adapter/primary/http/dto"
	"github.com/lechitz/aion-api/internal/user/core/domain"
	"github.com/lechitz/aion-api/internal/user/core/ports/input"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

const registrationIDParam = "registration_id"

// StartRegistration starts staged registration flow (step 1).
func (h *Handler) StartRegistration(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).Start(r.Context(), "user.registration.start")
	defer span.End()

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var req dto.StartRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	session, err := h.UserService.StartRegistration(ctx, input.StartRegistrationCommand{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "error starting registration", h.Logger)
		return
	}

	span.SetStatus(codes.Ok, "registration_started")
	httpresponse.WriteSuccess(w, http.StatusCreated, toRegistrationResponse(session), "registration started")
}

// UpdateRegistrationProfile updates staged registration step 2.
func (h *Handler) UpdateRegistrationProfile(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).Start(r.Context(), "user.registration.update_profile")
	defer span.End()

	registrationID := chi.URLParam(r, registrationIDParam)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req dto.UpdateRegistrationProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	session, err := h.UserService.UpdateRegistrationProfile(ctx, registrationID, input.UpdateRegistrationProfileCommand{
		Locale:   req.Locale,
		Timezone: req.Timezone,
		Location: req.Location,
		Bio:      req.Bio,
	})
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "error updating registration profile", h.Logger)
		return
	}

	span.SetStatus(codes.Ok, "registration_profile_updated")
	httpresponse.WriteSuccess(w, http.StatusOK, toRegistrationResponse(session), "registration profile updated")
}

// UpdateRegistrationAvatar updates staged registration step 3.
func (h *Handler) UpdateRegistrationAvatar(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).Start(r.Context(), "user.registration.update_avatar")
	defer span.End()

	registrationID := chi.URLParam(r, registrationIDParam)
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req dto.UpdateRegistrationAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteDecodeErrorSpan(ctx, w, span, err, h.Logger)
		return
	}

	session, err := h.UserService.UpdateRegistrationAvatar(ctx, registrationID, input.UpdateRegistrationAvatarCommand{
		AvatarURL: req.AvatarURL,
	})
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "error updating registration avatar", h.Logger)
		return
	}

	span.SetStatus(codes.Ok, "registration_avatar_updated")
	httpresponse.WriteSuccess(w, http.StatusOK, toRegistrationResponse(session), "registration avatar updated")
}

// CompleteRegistration completes staged registration and creates final user.
func (h *Handler) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
	ctx, span := otel.Tracer(TracerUserHandler).Start(r.Context(), "user.registration.complete")
	defer span.End()

	registrationID := chi.URLParam(r, registrationIDParam)
	user, err := h.UserService.CompleteRegistration(ctx, registrationID)
	if err != nil {
		httpresponse.WriteDomainErrorSpan(ctx, w, span, err, "error completing registration", h.Logger)
		return
	}

	span.SetStatus(codes.Ok, "registration_completed")
	httpresponse.WriteSuccess(w, http.StatusCreated, dto.CreateUserResponse{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		ID:       user.ID,
	}, "registration completed")
}

func toRegistrationResponse(session domain.RegistrationSession) dto.RegistrationSessionResponse {
	return dto.RegistrationSessionResponse{
		RegistrationID: session.RegistrationID,
		CurrentStep:    session.CurrentStep,
		Name:           session.Name,
		Username:       session.Username,
		Email:          session.Email,
		Locale:         session.Locale,
		Timezone:       session.Timezone,
		Location:       session.Location,
		Bio:            session.Bio,
		AvatarURL:      session.AvatarURL,
		ExpiresAt:      session.ExpiresAt.UTC().Format(timeLayoutRFC3339),
		Status:         string(session.Status),
	}
}

const timeLayoutRFC3339 = "2006-01-02T15:04:05Z07:00"
