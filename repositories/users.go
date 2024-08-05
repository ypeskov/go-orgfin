package repositories

import (
	"ypeskov/go-password-manager/internal/database"
	"ypeskov/go-password-manager/models"
)

type UsersRepository interface {
	GetAllUsers() ([]*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
	GetUserById(id int) (*models.User, error)
}

type usersRepoInstance struct {
	db *database.DbService
}

func NewUsersRepo(db *database.DbService) UsersRepository {
	return &usersRepoInstance{
		db: db,
	}
}

func (u *usersRepoInstance) CreateUser(user *models.User) error {
	log.Infof("Creating user: %+v\n", user)
	_, err := u.db.Db.NamedExec(`INSERT INTO users (email, name, hash_password) 
									VALUES (:email, :name, :hash_password)`, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *usersRepoInstance) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := u.db.Db.Select(&users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *usersRepoInstance) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := u.db.Db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *usersRepoInstance) GetUserById(userId int) (*models.User, error) {
	var user models.User

	err := u.db.Db.Get(&user, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
