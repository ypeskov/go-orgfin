package services

import (
	"ypeskov/go-orgfin/internal/database"
	"ypeskov/go-orgfin/internal/logger"
	"ypeskov/go-orgfin/repositories"
)

type ServiceManager struct {
	logger          *logger.Logger
	PasswordService PasswordsService
}

var managerInstance *ServiceManager

func NewServiceManager(db *database.Service, logger *logger.Logger) *ServiceManager {
	passwordRepo := repositories.NewPasswordRepo(db, logger)
	passwordService := NewPasswordService(&passwordRepo)

	managerInstance = &ServiceManager{
		logger:          logger,
		PasswordService: passwordService,
	}

	return managerInstance
}