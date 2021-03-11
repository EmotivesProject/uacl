package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/uacl/", func(r chi.Router) {
		r.Use(verifyJTW())
		r.Get("/", FetchItems)

		r.Get("/healthz", healthz)

		r.Route("/create_user", func(r chi.Router) {
			r.Post("/", CreateUser)
		})

		r.Route("/login", func(r chi.Router) {
			r.Post("/", Login)
		})

	})

	return r
}
