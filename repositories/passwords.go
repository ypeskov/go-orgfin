package repositories

import "ypeskov/go-orgfin/internal/database"

type PasswordsRepository interface {
	GetAllPasswords() ([]string, error)
}

type Passwords struct {
	db *database.Service
}

func NewPasswordRepo(db *database.Service) PasswordsRepository {
	return &Passwords{
		db: db,
	}
}

func (p *Passwords) GetAllPasswords() ([]string, error) {
	passwords := []string{"password1444444", "password2", "password3"}

	return passwords, nil
}
