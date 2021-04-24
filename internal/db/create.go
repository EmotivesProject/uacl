package db

import (
	"context"
	"time"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func CreateNewUser(user *model.User) error {
	logger.Info("Creating new user")

	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		context.Background(),
		"INSERT INTO users(name,username,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)",
		user.Name,
		user.Username,
		user.Password,
		time.Now(),
		time.Now(),
	)

	return err
}
