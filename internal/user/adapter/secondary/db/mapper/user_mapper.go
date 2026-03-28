// Package mapper contains functions to convert between domain and database models.
package mapper

import (
	"time"

	"github.com/lechitz/aion-api/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/aion-api/internal/user/core/domain"

	"gorm.io/gorm"
)

// UserFromDB converts a model.UserDB object into a domain.User object.
func UserFromDB(user model.UserDB) domain.User {
	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return domain.User{
		ID:                  user.ID,
		Name:                user.Name,
		Username:            user.Username,
		Email:               user.Email,
		Password:            user.Password,
		Locale:              user.Locale,
		Timezone:            user.Timezone,
		Location:            user.Location,
		Bio:                 user.Bio,
		AvatarURL:           user.AvatarURL,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		DeletedAt:           deletedAt,
	}
}

// UserToDB converts a domain.User object into a model.UserDB object for database storage.
// Roles are managed by /admin context and are not part of user domain.
func UserToDB(user domain.User) model.UserDB {
	var deleted gorm.DeletedAt
	if user.DeletedAt != nil {
		deleted.Time = *user.DeletedAt
		deleted.Valid = true
	}

	return model.UserDB{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		// Roles are NOT mapped - managed separately in user_roles table
		Locale:              user.Locale,
		Timezone:            user.Timezone,
		Location:            user.Location,
		Bio:                 user.Bio,
		AvatarURL:           user.AvatarURL,
		OnboardingCompleted: user.OnboardingCompleted,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
		DeletedAt:           deleted,
	}
}
