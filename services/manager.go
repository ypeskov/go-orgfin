package services

import (
	"ypeskov/go-orgfin/internal/database"
	"ypeskov/go-orgfin/repositories"
)

type ServiceManager struct {
	PasswordService PasswordsService
}

func NewServiceManager(db *database.Service) *ServiceManager {
	passwordRepo := repositories.NewPasswordRepo(db)
	passwordService := NewPasswordService(&passwordRepo)

	return &ServiceManager{
		PasswordService: passwordService,
	}
}
