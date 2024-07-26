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

type passRepoInstance struct {
	db *database.Service
}

var log *logger.Logger

func NewPasswordRepo(db *database.Service, logger *logger.Logger) PasswordsRepository {
	log = logger

	return &passRepoInstance{
		db: db,
	}
}

func (p *passRepoInstance) GetAllPasswords() ([]*models.Password, error) {
	var passwords []*models.Password
	err := p.db.Db.Select(&passwords, "SELECT * FROM passwords")
	if err != nil {
		log.Errorln(fmt.Sprintf("Error getting all passwords: %v", err))
		return nil, err
	}

	return passwords, nil
}

func (p *passRepoInstance) GetPasswordById(id string) (*models.Password, error) {
	var password models.Password
	err := p.db.Db.Get(&password, "SELECT * FROM passwords WHERE id = $1", id)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error getting password by id: %v", err))
		return nil, err
	}

	return &password, nil
}

func (p *passRepoInstance) AddPassword(password *models.Password) error {
	_, err := p.db.Db.NamedExec("INSERT INTO passwords (name, url, password) VALUES (:name, :url, :password)",
		password)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error adding password: %v", err))
		return err
	}

	return nil
}

func (p *passRepoInstance) UpdatePassword(password *models.Password) error {
	_, err := p.db.Db.NamedExec("UPDATE passwords SET name = :name, url = :url, password = :password WHERE id = :id",
		password)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error updating password: %v", err))
		return err
	}

	return nil
}
