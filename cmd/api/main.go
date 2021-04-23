package main

import (
	"log"
	"net/http"
	"os"

	"uacl/internal/api"
	"uacl/internal/db"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"

	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger("uacl")

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,POST,OPTIONS",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header",
	})

	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	db.ConnectDB()

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
