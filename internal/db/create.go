package db

import (
	"context"
	"time"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func CreateNewUser(ctx context.Context, user *model.User) error {
	logger.Info("Creating new user")

	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		ctx,
		"INSERT INTO users(name,username,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5)",
		user.Name,
		user.Username,
		user.Password,
		time.Now(),
		time.Now(),
	)

	return err
}

func CreateNewFollow(ctx context.Context, user, follow string) error {
	logger.Info("Creating new user")

	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		ctx,
		"INSERT INTO followers(username,follow_username) VALUES ($1,$2)",
		user,
		follow,
	)

	return err
}

func CreateNewAutologinToken(ctx context.Context, username, uuid string) error {
	logger.Info("Creating new autologin")

	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		ctx,
		"INSERT INTO autologin_tokens(username,autologin_token) VALUES ($1,$2)",
		username,
		uuid,
	)

	return err
}
