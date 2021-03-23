package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)

		r.Get("/public_key", publicKey)

		r.Get("/authorize", authorizeHeader)

		r.Route("/user", func(r chi.Router) {
			r.Post("/", createUser)

			r.Route("/{encoded_id}", func(r chi.Router) {
				r.Get("/", getUserByEncodedID)
			})
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", login)
		})

	})

	return r
}
