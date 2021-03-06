package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func CreateRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/", homeLink)
	return r
}
