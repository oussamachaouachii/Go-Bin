package mocks

import (
	"snippetbox.oussama.com/internal/models"
)

type UserModel struct{}

func (m *UserModel) Insert(name, email, password string) error {
	switch email {
	case "duplicate@app.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}
func (m *UserModel) Authenticate(email, password string) (int, string, error) {
	if email == "existing@app.com" && password == "password123" {
		return 1, "name", nil
	}
	return 0, "", models.ErrInvalidCredentials
}
func (m *UserModel) Exists(id int) (bool, error) {
	switch id {
	case 1:
		return true, nil
	default:
		return false, nil
	}
}



