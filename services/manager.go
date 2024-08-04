package services

import (
	"ypeskov/go-password-manager/internal/database"
	"ypeskov/go-password-manager/internal/logger"
	"ypeskov/go-password-manager/repositories"
)

type ServiceManager struct {
	logger          *logger.Logger
	PasswordService EncryptedPasswordsService
	UsersService    UsersService
}

var log *logger.Logger

func NewServiceManager(db *database.DbService, logger *logger.Logger) *ServiceManager {
	log = logger

	passwordRepo := repositories.NewPasswordRepo(db, logger)
	passwordService := NewPasswordService(&passwordRepo)

	usersRepo := repositories.NewUsersRepo(db)
	usersService := NewUsersService(&usersRepo)

	return &ServiceManager{
		logger:          logger,
		PasswordService: passwordService,
		UsersService:    usersService,
	}
}
