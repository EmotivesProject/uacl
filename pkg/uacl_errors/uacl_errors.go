package uacl_errors

import "errors"

var (
	ErrInvalidCredentials       = errors.New("Invalid login credentials. Please try again")
	ErrInvalidNameOrPassword    = errors.New("Must only contain letters, numbers, spaces or underscores")
	ErrInvalidEmailOrNameLength = errors.New("Must be between 3 and 100 characters")
	ErrInvalidPasswordLength    = errors.New("Must be between 6 and 100 characters")
	ErrInvalidEmail             = errors.New("Must be a valid email")
)
