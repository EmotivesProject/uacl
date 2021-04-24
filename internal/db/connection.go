package db

import (
	"context"
	"os"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/jackc/pgx/v4"

	"github.com/joho/godotenv"
)

var (
	db *pgx.Conn
)

//ConnectDB function: Make database connection
func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		logger.Error(err)
	}

	databaseURL := os.Getenv("databaseURL")

	//Define DB connection string and connect
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		logger.Error(err)
	}

	logger.Info("Successfully connected to the database")
	db = conn
}

func GetDB() *pgx.Conn {
	return db
}
