package test

import (
	"log"
	"net/http/httptest"
	"os"
	"uacl/internal/api"

	"github.com/EmotivesProject/common/logger"
	"github.com/EmotivesProject/common/middlewares"
	commonPostgres "github.com/EmotivesProject/common/postgres"
)

var TS *httptest.Server

func CreateStringAtLength(length int) string {
	var str string
	for i := 0; i < length; i++ {
		str += "x"
	}

	return str
}

func SetUpIntegrationTest() {
	logger.InitLogger("uacl", logger.EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Level:    os.Getenv("EMAIL_LEVEL"),
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST,OPTIONS,DELETE",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token",
	})

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/uacl_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	// These are set to be correct for internal/* tests
	os.Setenv("PRIVATE_KEY", "./../../../jwt/private.key")
	os.Setenv("PUBLIC_KEY", "./../../../jwt/public.key")

	os.Setenv("SECRET", "test")

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	con := commonPostgres.GetDatabase()
	con.Close()

	TS.Close()
}
