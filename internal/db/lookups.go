package db

import (
	"uacl/model"
	"uacl/pkg/uacl_errors"

	"gorm.io/gorm"
)

func FindByEmail(email string, database *gorm.DB) (model.User, error) {
	user := &model.User{}

	if err := database.Where("email = ?", email).First(user).Error; err != nil {
		return *user, uacl_errors.ErrInvalidCredentials
	}

	return *user, nil
}

func FindByEncodedID(encodedID string, database *gorm.DB) (model.User, error) {
	user := &model.User{}

	if err := database.Where("encoded_id = ?", encodedID).First(user).Error; err != nil {
		return *user, uacl_errors.ErrInvalidCredentials
	}

	return *user, nil
}
