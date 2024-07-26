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

var log *logger.Logger

func NewServiceManager(db *database.DbService, logger *logger.Logger) *ServiceManager {
	log = logger

	passwordRepo := repositories.NewPasswordRepo(db, logger)
	passwordService := NewPasswordService(&passwordRepo)

	return &ServiceManager{
		logger:          logger,
		PasswordService: passwordService,
	}
}
