package repositories

import (
	"fmt"
	"ypeskov/go-orgfin/internal/database"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/models"
)

type PasswordsRepository interface {
	GetAllPasswords() ([]*models.Password, error)
}

type Passwords struct {
	db     *database.Service
	logger *logger.Logger
}

func NewPasswordRepo(db *database.Service, logger *logger.Logger) PasswordsRepository {
	return &Passwords{
		db:     db,
		logger: logger,
	}
}

func (p *Passwords) GetAllPasswords() ([]*models.Password, error) {
	p.logger.Info("Getting all passwords repository")
	var passwords []*models.Password
	err := p.db.Db.Select(&passwords, "SELECT * FROM passwords")
	if err != nil {
		p.logger.Errorln(fmt.Sprintf("Error getting all passwords: %v", err))
		return nil, err
	}

	return passwords, nil
}
