package model

import (
	"regexp"
	"time"
	"uacl/messages"
)

const (
	nameField     = "Name"
	usernameField = "Username"
	passwordField = "Password"
	noField       = ""
)

var generalCharacters = regexp.MustCompile("^[A-Za-z0-9 _]+$")

// User struct declaration.
type User struct {
	ID        int       `json:"id,omitempty"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_time"`
	UpdatedAt time.Time `json:"updated_time"`
	Secret    string    `json:"secret"`
}

func (u User) ValidateCreate() (string, error) {
	if err := isNameValid(u.Name); err != nil {
		return nameField, err
	}

	if err := isUsernameValid(u.Username); err != nil {
		return usernameField, err
	}

	if err := isPasswordValid(u.Password); err != nil {
		return passwordField, err
	}

	if err := isSecretValid(u.Secret); err != nil {
		return noField, err
	}

	return noField, nil
}

func (u User) ValidateLogin() (string, error) {
	if err := isUsernameValid(u.Username); err != nil {
		return usernameField, err
	}

	if err := isPasswordValid(u.Password); err != nil {
		return passwordField, err
	}

	return noField, nil
}

func isNameValid(e string) error {
	if len(e) < 3 || len(e) > 100 {
		return messages.ErrInvalidUsernameOrNameLength
	}

	if !generalCharacters.MatchString(e) {
		return messages.ErrInvalidCharacter
	}

	return nil
}

func isUsernameValid(e string) error {
	if len(e) < 3 || len(e) > 100 {
		return messages.ErrInvalidUsernameOrNameLength
	}

	if !generalCharacters.MatchString(e) {
		return messages.ErrInvalidCharacter
	}

	return nil
}

func isPasswordValid(e string) error {
	if len(e) < 6 || len(e) > 100 {
		return messages.ErrInvalidPasswordLength
	}

	if !generalCharacters.MatchString(e) {
		return messages.ErrInvalidCharacter
	}

	return nil
}

func isSecretValid(e string) error {
	if e != "qutCreate" {
		return messages.ErrInvalidSecret
	}

	return nil
}
