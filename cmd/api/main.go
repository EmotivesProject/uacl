package main

import (
	"log"
	"net/http"

	"uacl/internal/api"
)

func main() {
	router := api.CreateRouter()
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
