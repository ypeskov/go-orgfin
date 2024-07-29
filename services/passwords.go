package services

import (
	"ypeskov/go-password-manager/models"
	"ypeskov/go-password-manager/repositories"
)

type PasswordsService interface {
	GetAllPasswords() ([]*models.Password, error)
	GetPasswordById(id int) (*models.Password, error)
	AddPassword(password *models.Password) error
	UpdatePassword(password *models.Password) error
	DeletePassword(id string) error
}

type passwordServiceInstance struct {
	PasswordRepo repositories.PasswordsRepository
}

func NewPasswordService(passwordRepo *repositories.PasswordsRepository) PasswordsService {
	return &passwordServiceInstance{
		PasswordRepo: *passwordRepo,
	}
}

func (p *passwordServiceInstance) GetAllPasswords() ([]*models.Password, error) {
	passwords, err := p.PasswordRepo.GetAllPasswords()
	if err != nil {
		log.Errorln("Error getting all passwords: %e\n", err)
		return nil, err
	}

	return passwords, nil
}

func (p *passwordServiceInstance) GetPasswordById(id int) (*models.Password, error) {
	password, err := p.PasswordRepo.GetPasswordById(id)
	if err != nil {
		log.Errorf("Error getting password by id: %e\n", err)
		return nil, err
	}

	return password, nil
}

func (p *passwordServiceInstance) AddPassword(password *models.Password) error {
	err := p.PasswordRepo.AddPassword(password)
	if err != nil {
		log.Errorf("Error adding password: %e\n", err)
		return err
	}

	return nil
}

func (p *passwordServiceInstance) UpdatePassword(password *models.Password) error {
	err := p.PasswordRepo.UpdatePassword(password)
	if err != nil {
		log.Errorf("Error updating password: %e\n", err)
		return err
	}

	return nil
}

func (p *passwordServiceInstance) DeletePassword(id string) error {
	err := p.PasswordRepo.DeletePassword(id)
	if err != nil {
		log.Errorf("Error deleting password: %e\n", err)
		return err
	}

	return nil
}
