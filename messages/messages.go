package messages

import "errors"

const (
	HealthResponse = "Health OK"
)

var (
	ErrInvalidCredentials          = errors.New("Invalid login credentials. Please try again")
	ErrInvalidCharacter            = errors.New("Must only contain letters, numbers, spaces or underscores")
	ErrInvalidUsernameOrNameLength = errors.New("Must be between 3 and 100 characters")
	ErrInvalidSecret               = errors.New("Incorrect secret value")
	ErrInvalidPasswordLength       = errors.New("Must be between 6 and 100 characters")
	ErrFailedDecoding              = errors.New("Failed during decoding request")
	ErrFailedCrypting              = errors.New("Failed during encrypting password")
	ErrUnauthorised                = errors.New("Unauthorized")
)
