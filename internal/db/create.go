package db

import (
	"context"
	"time"
	"uacl/model"

	"github.com/EmotivesProject/common/logger"
	commonPostgres "github.com/EmotivesProject/common/postgres"
)

func CreateNewUser(ctx context.Context, user *model.User) error {
	logger.Info("Creating new user")

	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		ctx,
		"INSERT INTO users(name,username,password,created_at,updated_at,user_group) VALUES ($1,$2,$3,$4,$5,$6)",
		user.Name,
		user.Username,
		user.Password,
		time.Now(),
		time.Now(),
		user.UserGroup,
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
