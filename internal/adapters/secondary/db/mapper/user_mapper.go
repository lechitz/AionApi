package mapper

import (
	"github.com/lechitz/AionApi/internal/adapters/secondary/db/model"
	"gorm.io/gorm"
	"time"

	"github.com/lechitz/AionApi/internal/core/domain"
)

func FromDB(u model.UserDB) domain.UserDomain {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return domain.UserDomain{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ToDB(u domain.UserDomain) model.UserDB {
	var deleted gorm.DeletedAt
	if u.DeletedAt != nil {
		deleted.Time = *u.DeletedAt
		deleted.Valid = true
	}

	return model.UserDB{
		ID:        u.ID,
		Name:      u.Name,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deleted,
	}
}
