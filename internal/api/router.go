package api

import (
	"net/http"

	"github.com/EmotivesProject/common/middlewares"
	"github.com/EmotivesProject/common/response"
	"github.com/go-chi/chi"
)

func CreateRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", response.Healthz)

		r.Get("/authorize", authorizeHeader)

		r.Post("/refresh", refreshToken)

		r.Post("/user", createUser)

		r.Post("/login", login)

		r.Route("/autologin", func(r chi.Router) {
			r.Post("/", createLoginToken)
			r.Post("/{token}", authoriseLoginToken)
			r.Get("/", getAutologinTokens)
			r.Get("/latest", getLatestAutologinTokens)
			r.Get("/{token_id}", getAutologinToken)
			r.Delete("/{token}", deleteAutologinToken)
		})
	})

	return r
}
