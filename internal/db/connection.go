package db

import (
	"fmt"
	"os"
	"uacl/model"

	"github.com/TomBowyerResearchProject/common/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

//ConnectDB function: Make database connection
func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal(err)
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	//Define DB connection string and connect
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	connectedDb, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	// Migrate the schema
	err = connectedDb.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Successfully connected to the database")
	db = connectedDb
}

func GetDB() *gorm.DB {
	return db
}
