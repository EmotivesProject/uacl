package db

import (
	"uacl/internal/uacl_errors"
	"uacl/model"

	"gorm.io/gorm"
)

func FindByUsername(username string, database *gorm.DB) (model.User, error) {
	user := &model.User{}

	if err := database.Where("username = ?", username).First(user).Error; err != nil {
		return *user, uacl_errors.ErrInvalidCredentials
	}

	return *user, nil
}
