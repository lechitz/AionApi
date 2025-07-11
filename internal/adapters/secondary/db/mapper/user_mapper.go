package mapper

import (
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"

	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"

	"gorm.io/gorm"
)

// UserFromDB converts a model.UserDB object into a domain.UserDomain object. It extracts and maps all user properties including timestamps.
func UserFromDB(user model.UserDB) domain.UserDomain {
	var deletedAt *time.Time
	if user.DeletedAt.Valid {
		deletedAt = &user.DeletedAt.Time
	}

	return domain.UserDomain{
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

// UserToDB converts a domain.UserDomain object into a model.UserDB object for database storage. It maps all relevant fields including timestamps.
func UserToDB(user domain.UserDomain) model.UserDB {
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
