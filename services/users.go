package services

import (
	"ypeskov/go-password-manager/models"
	"ypeskov/go-password-manager/repositories"
)

type UsersService interface {
	GetAllUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
}

type usersServiceInstance struct {
	UsersRepo repositories.UsersRepository
}

func NewUsersService(usersRepo *repositories.UsersRepository) UsersService {
	return &usersServiceInstance{
		UsersRepo: *usersRepo,
	}
}

func (u *usersServiceInstance) CreateUser(user *models.User) error {
	err := u.UsersRepo.CreateUser(user)
	if err != nil {
		log.Errorf("Error creating user: %e\n", err)
		return err
	}

	return nil
}

func (u *usersServiceInstance) GetAllUsers() ([]*models.User, error) {
	users, err := u.UsersRepo.GetAllUsers()
	if err != nil {
		log.Errorf("Error getting all users: %e\n", err)
		return nil, err
	}

	return users, nil
}

func (u *usersServiceInstance) GetUserByEmail(email string) (*models.User, error) {
	user, err := u.UsersRepo.GetUserByEmail(email)
	if err != nil {
		log.Errorf("Error getting user by email: %e\n", err)
		return nil, err
	}

	return user, nil
}
