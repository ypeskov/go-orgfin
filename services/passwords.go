package services

import (
	"ypeskov/go-orgfin/models"
	"ypeskov/go-orgfin/repositories"
)

type PasswordsService interface {
	GetAllPasswords() ([]*models.Password, error)
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
