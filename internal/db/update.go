package db

import (
	"context"
	"uacl/model"

	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func UpsertToken(ctx context.Context, token *model.Token) error {
	db := commonPostgres.GetDatabase()

	_, err := db.Exec(
		ctx,
		`INSERT INTO tokens(username,token,refresh_token,updated_at) VALUES ($1,$2,$3,$4)
		ON CONFLICT (username) DO UPDATE SET token = $2, refresh_token = $3, updated_at = $4`,
		token.Username,
		token.Token,
		token.RefreshToken,
		token.UpdatedAt,
	)

	return err
}
