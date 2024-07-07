package models

import "github.com/go-playground/validator"

type Password struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Url      string `json:"url"`
	Password string `json:"password"`
}

func (p *Password) Validate() error {
	validate := validator.New()

	return validate.Struct(p)
}
