package db

import (
	"errors"
	"fmt"
	"uacl/model"
	"uacl/pkg/auth"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errInvalidCredentials = errors.New("Invalid login credentials. Please try again")
)

func FindOne(email, password string, database *gorm.DB) (model.Token, error) {
	user := &model.User{}
	var token model.Token

	if err := database.Where("email = ?", email).First(user).Error; err != nil {
		return token, errInvalidCredentials
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(err.Error())
		return token, errInvalidCredentials
	}

	tokenString, err := auth.CreateToken(*user)
	if err != nil {
		return token, err
	}

	token.Token = tokenString
	return token, nil
}
