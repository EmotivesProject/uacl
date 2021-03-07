package db

import (
	"fmt"
	"log"
	"os"
	"uacl/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

	//Define DB connection string and connect
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	fmt.Println("Successfully connected to Database! ALL SYSTEMS ARE GO")
	return db
}
