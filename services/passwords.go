package services

import (
	"ypeskov/go-orgfin/models"
	"ypeskov/go-orgfin/repositories"
)

type PasswordsService interface {
	GetAllPasswords() ([]*models.Password, error)
	GetPasswordById(id string) (*models.Password, error)
}

type Passwords struct {
	PasswordRepo repositories.PasswordsRepository
}

func NewPasswordService(passwordRepo *repositories.PasswordsRepository) PasswordsService {
	return &Passwords{
		PasswordRepo: *passwordRepo,
	}
}

func (p *Passwords) GetAllPasswords() ([]*models.Password, error) {
	managerInstance.logger.Info("Getting all passwords service")
	passwords, err := p.PasswordRepo.GetAllPasswords()
	if err != nil {
		return nil, err
	}

	return passwords, nil
}

func (p *Passwords) GetPasswordById(id string) (*models.Password, error) {
	managerInstance.logger.Info("Getting password by id service")
	password, err := p.PasswordRepo.GetPasswordById(id)
	if err != nil {
		managerInstance.logger.Errorf("Error getting password by id: %e\n", err)
		return nil, err
	}

	return password, nil
}
