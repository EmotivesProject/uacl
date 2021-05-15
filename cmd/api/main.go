package main

import (
	"log"
	"net/http"
	"os"
	"uacl/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

func main() {
	logger.InitLogger("uacl")

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,POST,OPTIONS",
		AllowedHeaders: "*",
	})

	router := api.CreateRouter()

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: os.Getenv("DATABASE_URL"),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Fatal(http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), router))
}
