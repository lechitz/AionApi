package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"strings"
	"time"
)

type User struct {
	ID        uint64    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (user *User) Prepare(step string) error {
	if err := user.validate(step); err != nil {
		return err
	}

	user.format()

	return nil
}

func (user *User) validate(step string) error {
	if user.Name == "" {
		return errors.New("Name is required")
	}
	if user.Username == "" {
		return errors.New("Username is required")
	}
	if user.Email == "" {
		return errors.New("Email is required")
	}

	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("Invalid email")
	}

	if step == "register" && user.Password == "" {
		return errors.New("Password is required")
	}
	return nil
}

func (user *User) format() {
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
}
