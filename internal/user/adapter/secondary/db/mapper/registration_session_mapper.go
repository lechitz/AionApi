package mapper

import (
	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"
)

// RegistrationSessionFromDB maps db model to domain object.
func RegistrationSessionFromDB(session model.RegistrationSessionDB) domain.RegistrationSession {
	return domain.RegistrationSession{
		RegistrationID: session.RegistrationID,
		Name:           session.Name,
		Username:       session.Username,
		Email:          session.Email,
		PasswordHash:   session.PasswordHash,
		Locale:         session.Locale,
		Timezone:       session.Timezone,
		Location:       session.Location,
		Bio:            session.Bio,
		AvatarURL:      session.AvatarURL,
		CurrentStep:    session.CurrentStep,
		Status:         domain.RegistrationSessionStatus(session.Status),
		ExpiresAt:      session.ExpiresAt,
		CreatedAt:      session.CreatedAt,
		UpdatedAt:      session.UpdatedAt,
	}
}

// RegistrationSessionToDB maps domain object to db model.
func RegistrationSessionToDB(session domain.RegistrationSession) model.RegistrationSessionDB {
	return model.RegistrationSessionDB{
		RegistrationID: session.RegistrationID,
		Name:           session.Name,
		Username:       session.Username,
		Email:          session.Email,
		PasswordHash:   session.PasswordHash,
		Locale:         session.Locale,
		Timezone:       session.Timezone,
		Location:       session.Location,
		Bio:            session.Bio,
		AvatarURL:      session.AvatarURL,
		CurrentStep:    session.CurrentStep,
		Status:         string(session.Status),
		ExpiresAt:      session.ExpiresAt,
		CreatedAt:      session.CreatedAt,
		UpdatedAt:      session.UpdatedAt,
	}
}
