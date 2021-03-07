package api

import (
	"github.com/go-chi/chi"
)

// can also do r.Use(authorizationMiddleware()) and r.With(authorization)
func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.With(SimpleMiddleware())

	r.Get("/healthz", healthz)

	r.Route("/create_user", func(r chi.Router) {
		r.Post("/", CreateUser)
	})

	r.Route("/login", func(r chi.Router) {
		r.Post("/", Login)
	})

	r.Route("/", func(r chi.Router) {
		r.Use(verifyJTW())
		r.Get("/", FetchItems)
	})

	return r
}
