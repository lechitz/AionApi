package mapper

import (
	"github.com/lechitz/AionApi/adapters/secondary/db/model"
	"gorm.io/gorm"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
)

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
