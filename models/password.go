package models

import "github.com/go-playground/validator"

type Password struct {
	Id       int    `json:"id" form:"id"`
	Name     string `json:"name" form:"name" validate:"required"`
	Url      string `json:"url" form:"url"`
	Password string `json:"password" form:"password" validate:"required"`
}

func (p *Password) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}
