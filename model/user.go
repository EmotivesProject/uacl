package model

import (
	"regexp"
	"time"
	"uacl/pkg/uacl_errors"
)

const (
	nameField     = "Name"
	emailField    = "Email"
	passwordField = "Password"
	noField       = ""
)

var (
	generalCharacters = regexp.MustCompile("[A-Za-z0-9 _]")
	emailRegex        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

//User struct declaration
type User struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_time" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_time" gorm:"autoUpdateTime"`
	EncodedID string    `json:"encoded_id"`
}

func (u User) ValidateCreate() (string, error) {
	if err := isNameValid(u.Name); err != nil {
		return nameField, err
	}
	if err := isEmailValid(u.Email); err != nil {
		return emailField, err
	}
	if err := isPasswordValid(u.Password); err != nil {
		return passwordField, err
	}
	return noField, nil
}

func (u User) ValidateLogin() (string, error) {
	if err := isEmailValid(u.Email); err != nil {
		return emailField, err
	}
	if err := isPasswordValid(u.Password); err != nil {
		return passwordField, err
	}
	return noField, nil
}

func isNameValid(e string) error {
	if len(e) < 3 && len(e) > 100 {
		return uacl_errors.ErrInvalidEmailOrNameLength
	}
	if !generalCharacters.MatchString(e) {
		return uacl_errors.ErrInvalidNameOrPassword
	}
	return nil
}

func isEmailValid(e string) error {
	if len(e) < 3 && len(e) > 100 {
		return uacl_errors.ErrInvalidEmailOrNameLength
	}
	if !emailRegex.MatchString(e) {
		return uacl_errors.ErrInvalidEmail
	}
	return nil
}

func isPasswordValid(e string) error {
	if len(e) < 3 && len(e) > 100 {
		return uacl_errors.ErrInvalidPasswordLength
	}
	if !generalCharacters.MatchString(e) {
		return uacl_errors.ErrInvalidNameOrPassword
	}
	return nil
}
