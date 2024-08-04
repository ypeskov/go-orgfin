package models

import (
	"github.com/go-playground/validator"
	"time"
)

type User struct {
	Id           int        `json:"id" form:"id" db:"id"`
	Name         string     `json:"name" form:"name" db:"name"`
	Email        string     `json:"email" form:"email" db:"email" validate:"required,email"`
	HashPassword string     `json:"password" form:"password" db:"hash_password" validate:"required"`
	LastLogin    *time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

func (u *User) Validate() error {
	validate := validator.New()

	return validate.Struct(u)
}

// We use HTMX, so we don't need to return a JSON response or process AJAX requests

//type UserRequest struct {
//	Name     string `json:"name" form:"name" validate:"required"`
//	Email    string `json:"email" form:"email" validate:"required,email"`
//	EncryptedPassword string `json:"password" form:"password" validate:"required"`
//}
//
//func (u *UserRequest) Validate() error {
//	validate := validator.New()
//	return validate.Struct(u)
//}
//
//type UserResponse struct {
//	Id    int    `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
//
//func ToUserModel(userRequest *UserRequest) *User {
//	return &User{
//		Name:     userRequest.Name,
//		Email:    userRequest.Email,
//		EncryptedPassword: userRequest.EncryptedPassword,
//	}
//}
//
//func ToUserResponse(user *User) *UserResponse {
//	return &UserResponse{
//		Id:    user.Id,
//		Name:  user.Name,
//		Email: user.Email,
//	}
//}
