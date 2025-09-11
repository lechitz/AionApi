package mapper

import (
	"time"

	"github.com/lechitz/AionApi/internal/user/adapter/secondary/db/model"
	"github.com/lechitz/AionApi/internal/user/core/domain"

	"gorm.io/gorm"
)

// UserFromDB converts a model.UserDB object into a domain.User object. It extracts and maps all user properties including timestamps.
func UserFromDB(user model.UserDB) domain.User {
	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

// UserToDB converts a domain.User object into a model.UserDB object for database storage. It maps all relevant fields including timestamps.
func UserToDB(user domain.User) model.UserDB {
	var deleted gorm.DeletedAt
	if user.DeletedAt != nil {
		deleted.Time = *user.DeletedAt
		deleted.Valid = true
	}

	return model.UserDB{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deleted,
	}
}
