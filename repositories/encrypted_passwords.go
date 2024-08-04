package repositories

import (
	"fmt"
	"ypeskov/go-password-manager/internal/database"
	"ypeskov/go-password-manager/internal/logger"
	"ypeskov/go-password-manager/models"
)

type EncryptedPasswordsRepository interface {
	GetAllPasswords(userId int) ([]*models.EncryptedPassword, error)
	GetPasswordById(id int) (*models.EncryptedPassword, error)
	AddPassword(password *models.EncryptedPassword) error
	UpdatePassword(password *models.EncryptedPassword) error
	DeletePassword(id string) error
}

type passRepoInstance struct {
	db *database.DbService
}

var log *logger.Logger

func NewPasswordRepo(db *database.DbService, logger *logger.Logger) EncryptedPasswordsRepository {
	log = logger

	return &passRepoInstance{
		db: db,
	}
}

func (p *passRepoInstance) GetAllPasswords(userId int) ([]*models.EncryptedPassword, error) {
	var passwords []*models.EncryptedPassword
	err := p.db.Db.Select(&passwords, "SELECT * FROM encrypted_passwords WHERE user_id = $1", userId)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error getting all passwords: %v", err))
		return nil, err
	}

	return passwords, nil
}

func (p *passRepoInstance) GetPasswordById(id int) (*models.EncryptedPassword, error) {
	var password models.EncryptedPassword

	err := p.db.Db.Get(&password, "SELECT * FROM encrypted_passwords WHERE id = $1", id)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error getting password by id: %v", err))
		return nil, err
	}

	return &password, nil
}

func (p *passRepoInstance) AddPassword(password *models.EncryptedPassword) error {
	log.Infof("Adding password for user id: %v\n", password.UserId)
	_, err := p.db.Db.NamedExec(`INSERT INTO encrypted_passwords (user_id, name, resource, password, salt, iv) 
		VALUES (:user_id, :name, :resource, :password, :salt, :iv)`,
		password)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error adding password: %v", err))
		return err
	}

	return nil
}

func (p *passRepoInstance) UpdatePassword(password *models.EncryptedPassword) error {
	_, err := p.db.Db.NamedExec(`UPDATE passwords SET name = :name, resource = :resource,
                     password = :password, salt = :salt, iv = :iv WHERE id = :id`,
		password)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error updating password: %v", err))
		return err
	}

	return nil
}

func (p *passRepoInstance) DeletePassword(id string) error {
	_, err := p.db.Db.Exec("DELETE FROM passwords WHERE id = $1", id)
	if err != nil {
		log.Errorln(fmt.Sprintf("Error deleting password: %v", err))
		return err
	}

	return nil
}
