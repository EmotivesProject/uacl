package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(givenPassword, databasePassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(givenPassword))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	return true
}

func CreatePassword(password string) string {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(encryptedPassword)
}
