package db

import (
	"context"
	"uacl/model"
)

func FindByUsername(username string) (model.User, error) {
	user := model.User{}
	db = GetDB()

	err := db.QueryRow(context.Background(), "SELECT id,name,username,password,created_at,updated_at FROM users WHERE username = $1", username).Scan(
		&user.ID, &user.Name, &user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)

	return user, err
}
