package uacl_errors

import "errors"

var (
	ErrInvalidCredentials          = errors.New("Invalid login credentials. Please try again")
	ErrInvalidCharacter            = errors.New("Must only contain letters, numbers, spaces or underscores")
	ErrInvalidUsernameOrNameLength = errors.New("Must be between 3 and 100 characters")
	ErrInvalidSecret               = errors.New("Incorrect secret value")
	ErrInvalidPasswordLength       = errors.New("Must be between 6 and 100 characters")
)
