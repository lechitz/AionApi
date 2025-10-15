// Package mapper contains functions to convert between domain and database models.
package mapper

import (
	"strings"
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

	// Convert comma-separated roles string to slice
	var roles []string
	if user.Roles != "" {
		roles = strings.Split(user.Roles, ",")
		// Trim spaces from each role
		for i := range roles {
			roles[i] = strings.TrimSpace(roles[i])
		}
	}

	return domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     roles,
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

	// Convert roles slice to comma-separated string
	rolesStr := "user" // default value
	if len(user.Roles) > 0 {
		rolesStr = strings.Join(user.Roles, ",")
	}

	return model.UserDB{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		Roles:     rolesStr,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: deleted,
	}
}
