package repositories

import (
	"fmt"
	"ypeskov/go-orgfin/internal/database"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/models"
)

type PasswordsRepository interface {
	GetAllPasswords() ([]*models.Password, error)
	GetPasswordById(id string) (*models.Password, error)
	AddPassword(password *models.Password) error
	UpdatePassword(password *models.Password) error
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

func (p *Passwords) GetPasswordById(id string) (*models.Password, error) {
	p.logger.Info("Getting password by id repository")
	var password models.Password
	err := p.db.Db.Get(&password, "SELECT * FROM passwords WHERE id = $1", id)
	if err != nil {
		p.logger.Errorln(fmt.Sprintf("Error getting password by id: %v", err))
		return nil, err
	}

	return &password, nil
}

func (p *Passwords) AddPassword(password *models.Password) error {
	p.logger.Info("Adding password repository")

	_, err := p.db.Db.NamedExec("INSERT INTO passwords (name, url, password) VALUES (:name, :url, :password)",
		password)
	if err != nil {
		p.logger.Errorln(fmt.Sprintf("Error adding password: %v", err))
		return err
	}

	return nil
}

func (p *Passwords) UpdatePassword(password *models.Password) error {
	p.logger.Info("Updating password repository")

	_, err := p.db.Db.NamedExec("UPDATE passwords SET name = :name, url = :url, password = :password WHERE id = :id",
		password)
	if err != nil {
		p.logger.Errorln(fmt.Sprintf("Error updating password: %v", err))
		return err
	}

	return nil
}
