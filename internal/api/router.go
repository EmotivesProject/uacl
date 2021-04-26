package api

import (
	"github.com/TomBowyerResearchProject/common/middlewares"
	"github.com/go-chi/chi"
)

func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middlewares.SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)

		r.Get("/public_key", publicKey)

		r.Get("/authorize", authorizeHeader)

		r.Route("/user", func(r chi.Router) {
			r.Post("/", createUser)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", login)
		})
	})

	return r
}
