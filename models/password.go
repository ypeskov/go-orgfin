package models

import "github.com/go-playground/validator"

type Password struct {
	Id       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name" validate:"required"`
	Resource string `json:"resource" form:"resource"`
	Password string `json:"password" form:"password" validate:"required"`
	Salt     string `json:"salt" form:"salt"`
	Iv       string `json:"iv" form:"iv"`
}

func (p *Password) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}
