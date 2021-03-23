package db

import (
	"uacl/model"
	"uacl/pkg/uacl_errors"

	"gorm.io/gorm"
)

func FindByUsername(username string, database *gorm.DB) (model.User, error) {
	user := &model.User{}

	if err := database.Where("username = ?", username).First(user).Error; err != nil {
		return *user, uacl_errors.ErrInvalidCredentials
	}

	return *user, nil
}
