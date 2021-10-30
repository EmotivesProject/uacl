package db

import (
	"context"
	"errors"
	"os"
	"uacl/messages"
	"uacl/model"

	"github.com/EmotivesProject/common/logger"
	commonPostgres "github.com/EmotivesProject/common/postgres"
	"github.com/jackc/pgx/v4"
)

func FindByUsername(ctx context.Context, username string) (model.User, error) {
	user := model.User{}
	db := commonPostgres.GetDatabase()

	err := db.QueryRow(
		ctx,
		"SELECT id,name,username,password,created_at,updated_at,user_group FROM users WHERE username = $1",
		username,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.UserGroup,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return user, messages.ErrInvalidCredentials
	}

	return user, err
}

func RefreshTokenIsValidForUsername(ctx context.Context, refreshToken, username string) bool {
	db := commonPostgres.GetDatabase()

	var exists bool

	err := db.QueryRow(
		ctx,
		"SELECT exists (SELECT * FROM tokens where username = $1 and refresh_token = $2)",
		username,
		refreshToken,
	).Scan(&exists)
	if err != nil {
		logger.Error(err)

		return false
	}

	return exists
}

func FindAutologinForUser(ctx context.Context, autologinToken string) (model.AutologinRequest, error) {
	db := commonPostgres.GetDatabase()

	var user model.AutologinRequest

	err := db.QueryRow(
		ctx,
		"SELECT username FROM autologin_tokens where autologin_token = $1",
		autologinToken,
	).Scan(&user.Username)

	return user, err
}

func FindAutologinByID(ctx context.Context, autologinID int) (model.AutologinToken, error) {
	db := commonPostgres.GetDatabase()

	var autologin model.AutologinToken

	err := db.QueryRow(
		ctx,
		"SELECT username, autologin_token FROM autologin_tokens where id = $1",
		autologinID,
	).Scan(&autologin.Username, &autologin.AutologinToken)

	autologin.Site = os.Getenv("AUTOLOGIN_URL")

	return autologin, err
}

func FindAutologins(ctx context.Context) ([]model.AutologinToken, error) {
	db := commonPostgres.GetDatabase()

	autologins := make([]model.AutologinToken, 0)

	rows, err := db.Query(
		ctx,
		"SELECT username, autologin_token FROM autologin_tokens",
	)
	if err != nil {
		return autologins, err
	}

	for rows.Next() {
		var autologin model.AutologinToken

		err := rows.Scan(
			&autologin.Username,
			&autologin.AutologinToken,
		)
		if err != nil {
			continue
		}

		autologin.Site = os.Getenv("AUTOLOGIN_URL")

		autologins = append(autologins, autologin)
	}

	return autologins, err
}
