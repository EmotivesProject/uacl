package api

import (
	"context"
	"errors"
	"net/http"
	"uacl/model"
	"uacl/pkg/auth"
)

type key string

const (
	userID key = "username"
)

var (
	errUnauthorised = errors.New("Not authorised")
)

func verifyJTW() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			user, err := auth.Validate(header)
			if err != nil {
				messageResponseJSON(w, http.StatusBadRequest, model.Message{
					Message: errUnauthorised.Error(),
				})
				return
			}
			ctx := context.WithValue(r.Context(), userID, user.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func SimpleMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			w.Header().Set("Access-Control-Allow-Origin", "*") // fixme please
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")

			if r.Method == "OPTIONS" {
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
