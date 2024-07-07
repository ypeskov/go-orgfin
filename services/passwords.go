package services

import "ypeskov/go-orgfin/repositories"

type PasswordsService interface {
	GetAllPasswords() ([]string, error)
}

type Passwords struct {
	PasswordRepo repositories.PasswordsRepository
}

func NewPasswordService(passwordRepo *repositories.PasswordsRepository) PasswordsService {
	return &Passwords{
		PasswordRepo: *passwordRepo,
	}
}

func (p *Passwords) GetAllPasswords() ([]string, error) {
	passwords, err := p.PasswordRepo.GetAllPasswords()
	if err != nil {
		return nil, err
	}

	return passwords, nil
}
