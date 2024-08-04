package models

import (
	"github.com/go-playground/validator"
	"time"
)

type EncryptedPassword struct {
	Id        int       `json:"id" form:"id" db:"id"`
	UserId    int       `json:"user_id" form:"user_id" db:"user_id"`
	Name      string    `json:"name" form:"name" validate:"required" db:"name"`
	Login     string    `json:"login" form:"login" db:"login"`
	Resource  string    `json:"resource" form:"resource" db:"resource"`
	Password  string    `json:"password" form:"password" validate:"required" db:"password"`
	Salt      string    `json:"salt" form:"salt" db:"salt"`
	Iv        string    `json:"iv" form:"iv" db:"iv"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

func (p *EncryptedPassword) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}
