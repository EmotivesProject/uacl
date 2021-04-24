package main

import (
	"log"
	"net/http"
	"os"

	"uacl/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"

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

	commonPostgres.Connect(commonPostgres.Config{
		URI: os.Getenv("databaseURL"),
	})

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
