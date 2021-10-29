package db

import (
	"context"

	commonPostgres "github.com/EmotivesProject/common/postgres"
)

func DeleteAutologinToken(ctx context.Context, token string) error {
	db := commonPostgres.GetDatabase()

	_, err := db.Query(
		ctx,
		"DELETE FROM autologin_tokens where autologin_token = $1",
		token,
	)

	return err
}
