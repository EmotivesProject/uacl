package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/uacl/", func(r chi.Router) {
		r.Get("/healthz", healthz)

		r.Get("/public_key", publicKey)

		r.Route("/user", func(r chi.Router) {
			r.Post("/", CreateUser)
			// r.Put("/", UpdateUser)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", Login)
		})

	})

	return r
}
