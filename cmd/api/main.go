package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uacl/internal/api"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
)

const timeBeforeTimeout = 15

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

	srv := http.Server{
		Handler:      router,
		Addr:         os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		WriteTimeout: timeBeforeTimeout * time.Second,
		ReadTimeout:  timeBeforeTimeout * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		logger.Info("Graceful shutdown initiated")

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Infof("HTTP server shutdown: %v", err)
		}

		commonPostgres.CloseDatabase()
		logger.Info("Shutdown the server and disconnected postgres")

		close(idleConnsClosed)
	}()

	logger.Infof("Listening for http on: %s", os.Getenv("HOST")+":"+os.Getenv("PORT"))

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		commonPostgres.CloseDatabase()
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
