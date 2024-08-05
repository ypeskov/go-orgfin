package errors

type UserError struct {
	Message string
}

func (e *UserError) UserError() string {
	return e.Message
}

var UserNotFound = &UserError{Message: "User not found"}
var UserExists = &UserError{Message: "User already exists"}
var UserPasswordsDoNotMatch = &UserError{Message: "Passwords do not match"}
var MissingEmail = &UserError{Message: "Missing email"}
var MissingPassword = &UserError{Message: "Missing password"}
var MissingConfirmPassword = &UserError{Message: "Missing confirm password"}
