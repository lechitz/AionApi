package repository

import (
	"database/sql"
	"fmt"
	"github.com/lechitz/AionApi/src/models"
)

type Users struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *Users {
	return &Users{db: db}
}

func (repository Users) CreateUser(user models.User) (uint64, error) {
	statement, err := repository.db.Prepare("insert into users (name, username, password, email) values (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Username, user.Password, user.Email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func (repository Users) GetUser(userOrName string) ([]models.User, error) {
	nickOrName := fmt.Sprintf("%%%s%%", userOrName)

	lines, err := repository.db.Query(
		"select id, name, username, email, created_at from users where name like ? or username like ?",
		nickOrName,
		nickOrName,
	)

	if err != nil {
		return nil, err
	}
	defer lines.Close()

	var users []models.User

	for lines.Next() {
		var user models.User

		if err = lines.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository Users) GetUserByID(id uint64) (models.User, error) {
	lines, err := repository.db.Query("select id, name, username, email, created_at from users where id = ?", id)
	if err != nil {
		return models.User{}, err
	}
	defer lines.Close()

	var user models.User

	if lines.Next() {
		if err = lines.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.CreatedAt); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repository Users) UpdateUser(ID uint64, user models.User) error {
	statement, err := repository.db.Prepare("update users set name = ?, username = ?, email = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(user.Name, user.Username, user.Email, ID); err != nil {
		return err
	}

	return nil
}

func (repository Users) DeleteUser(ID uint64) error {
	statement, err := repository.db.Prepare("delete from users where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(ID); err != nil {
		return err
	}

	return nil
}

func (repository Users) GetUserByEmail(email string) (models.User, error) {
	lines, err := repository.db.Query("select id, password from users where email = ?", email)
	if err != nil {
		return models.User{}, err
	}
	defer lines.Close()

	var user models.User

	if lines.Next() {
		if err = lines.Scan(&user.ID, &user.Password); err != nil {
			return models.User{}, err
		}
	}

	return user, nil
}

func (repository Users) GetPassword(userID uint64) (string, error) {
	line, err := repository.db.Query("select password from users where id = ?", userID)
	if err != nil {
		return "", err
	}
	defer line.Close()

	var user models.User

	if line.Next() {
		if err = line.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (repository Users) UpdatePassword(userID uint64, newPassword string) error {
	statement, err := repository.db.Prepare("update users set password = ? where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err = statement.Exec(newPassword, userID); err != nil {
		return err
	}

	return nil
}
