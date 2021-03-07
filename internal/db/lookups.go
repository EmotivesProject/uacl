package db

import (
	"errors"
	"fmt"
	"time"
	"uacl/model"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	errNoEmail            = errors.New("No email found in request")
	errInvalidCredentials = errors.New("Invalid login credentials. Please try again")
)

func FindOne(email, password string, database *gorm.DB) (map[string]interface{}, error) {
	user := &model.User{}

	if err := database.Where("email = ?", email).First(user).Error; err != nil {
		return nil, errNoEmail
	}

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, errInvalidCredentials
	}

	fmt.Println(err.Error())

	tk := &model.Token{
		EncodedID: user.EncodedID,
		Name:      user.Name,
		Email:     user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}

	response := map[string]interface{}{
		"token": tokenString,
	}
	return response, nil
}
