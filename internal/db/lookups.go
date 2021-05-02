package db

import (
	"context"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func FindByUsername(username string) (model.User, error) {
	user := model.User{}
	db := commonPostgres.GetDatabase()

	err := db.QueryRow(
		context.Background(),
		"SELECT id,name,username,password,created_at,updated_at FROM users WHERE username = $1",
		username,
	).Scan(
		&user.ID, &user.Name, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}

func RefreshTokenIsValidForUsername(refreshToken, username string) bool {
	db := commonPostgres.GetDatabase()

	var exists bool

	err := db.QueryRow(
		context.TODO(),
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
