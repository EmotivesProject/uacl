package main

import (
	"log"
	"net/http"
	"os"

	"uacl/internal/api"
	"uacl/internal/db"
	"uacl/internal/logger"

	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger()

	db.ConnectDB()

	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
