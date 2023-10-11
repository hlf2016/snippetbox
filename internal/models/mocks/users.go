package mocks

import (
	"github.com/hlf2016/snippetbox/internal/models"
	"time"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "dupe@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, error) {
	if email == "alice@example.com" && password == "password" {
		return 1, nil
	}
	return 0, models.ErrInvalidCredential
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}

func (m *UserModel) Get(id int) (*models.User, error) {
	if id == 1 {
		u := &models.User{
			ID:      1,
			Name:    "test",
			Email:   "example@email.com",
			Created: time.Now(),
		}
		return u, nil
	}
	return nil, models.ErrNoRecord
}

func (m *UserModel) PasswordUpdate(id int, currentPassword, newPassword string) error {
	if id == 1 {
		if currentPassword == "password" {
			return models.ErrInvalidCredential
		}
		return nil
	}

	return models.ErrNoRecord
}
