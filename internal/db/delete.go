package db

import (
	"context"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
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

func DeleteFollowByUserAndFollow(ctx context.Context, user, follow string) error {
	db := commonPostgres.GetDatabase()

	_, err := db.Query(
		ctx,
		"DELETE FROM followers where username = $1 and follow_username = $2",
		user,
		follow,
	)

	return err
}
